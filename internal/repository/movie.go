package repository

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
	"truckx/internal/models"

	"truckx/internal/services/util"

	"github.com/jmoiron/sqlx"
)

type MovieRepo interface {
	Create(ctx context.Context, movie models.Movie) (*models.Movie, error)
	List(ctx context.Context) ([]models.Movie, error)
	GetAvailability(ctx context.Context, movieId int, state string, date time.Time) ([]models.Availability, error)
}

func NewMovieRepo(db *sqlx.DB) MovieRepo {
	return &movieImpl{db}
}

type movieImpl struct {
	db *sqlx.DB
}

func (r movieImpl) Create(ctx context.Context, movie models.Movie) (*models.Movie, error) {
	row, err := r.db.QueryContext(ctx,
		"INSERT INTO movies (name, release_date) VALUES ($1, $2) RETURNING id",
		movie.Name, movie.ReleaseDate)
	if err != nil {
		return nil, err
	}
	err = row.Scan(&movie.Id)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r movieImpl) Get(ctx context.Context, id int) (*models.Movie, error) {
	movie := models.Movie{}
	row := r.db.QueryRowxContext(ctx, "SELECT * FROM movies WHERE id = $1", id)
	err := row.StructScan(&movie)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r movieImpl) List(ctx context.Context) ([]models.Movie, error) {
	movies := []models.Movie{}
	row, err := r.db.QueryxContext(ctx, "SELECT * FROM movies")
	if err != nil {
		return nil, err
	}
	var movie models.Movie
	for row.Next() {
		err = row.StructScan(&movie)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (r movieImpl) GetAvailability(ctx context.Context, movieId int, state string, date time.Time) ([]models.Availability, error) {

	theatreStatQuery := `SELECT ms.id movie_show_id, theatre_details.* FROM movies_shows ms JOIN 
	(
	SELECT t.id theatre_id, t.name theatre_name, s.id show_id, s.from_hr show_from_hr, s.from_min show_from_min, 
	s.till_hr show_till_hr, s.till_min show_till_min, 
	array_to_json(array_agg(st.id)::int[]) theatre_seat_ids FROM theatres t JOIN shows s ON t.id = s.theatre_id 
	JOIN seats st ON t.id = st.theatre_id WHERE t.state = $1 GROUP BY t.id, s.id
	) theatre_details 
 ON ms.show_id = theatre_details.show_id
WHERE ms.movie_id = $2`

	row, err := r.db.QueryxContext(ctx, theatreStatQuery, state, movieId)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var s stat
	theatreStats := make(map[int][]stat)
	for row.Next() {
		err = row.StructScan(&s)
		if err != nil {
			return nil, err
		}
		if stats, ok := theatreStats[s.TheatreId]; ok {
			stats = append(stats, s)
			theatreStats[s.TheatreId] = stats
		} else {
			theatreStats[s.TheatreId] = []stat{s}
		}
	}

	bookingStatQuery := `SELECT bs.seat_id booked_seat_id, booking_theatre_show_details.* FROM bookings_seats bs JOIN (
		SELECT b.id as booking_id, theatre_show_details.* FROM bookings b 
		JOIN 
		   (SELECT ms.id movie_show_id, theatre_details.* FROM movies_shows ms JOIN 
			  (
			  SELECT t.id theatre_id, t.name theatre_name, s.id show_id, s.from_hr show_from_hr, s.from_min show_from_min, 
			  s.till_hr show_till_hr, s.till_min show_till_min, 
			  array_to_json(array_agg(st.id)::int[]) theatre_seat_ids FROM theatres t JOIN shows s ON t.id = s.theatre_id 
			  JOIN seats st ON t.id = st.theatre_id WHERE t.state = $1 GROUP BY t.id, s.id
			  ) theatre_details 
		   ON ms.show_id = theatre_details.show_id
		WHERE ms.movie_id = $2) theatre_show_details 
		ON b.movies_shows_id = theatre_show_details.movie_show_id
		WHERE b.booking_date::date = $3
  ) booking_theatre_show_details 
  ON bs.booking_id = booking_theatre_show_details.booking_id;`

	row, err = r.db.QueryxContext(ctx, bookingStatQuery, state, movieId, time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()))
	if err != nil {
		return nil, err
	}
	defer row.Close()
	bookingStats := make(map[int][]stat)
	for row.Next() {
		err = row.StructScan(&s)
		if err != nil {
			return nil, err
		}
		if stats, ok := bookingStats[s.TheatreId]; ok {
			stats = append(stats, s)
			bookingStats[s.TheatreId] = stats
		} else {
			bookingStats[s.TheatreId] = []stat{s}
		}
	}

	theatreWiseAvalibilities := make(map[int]models.Availability)
	for theatreId, tstats := range theatreStats {
		if bStats, ok := bookingStats[theatreId]; ok {
			for _, theatreStat := range tstats {
				slot := fmt.Sprintf("%d:%d-%d:%d", theatreStat.ShowFromHr, theatreStat.ShowFromMin, theatreStat.ShowTillHr, theatreStat.ShowTillMin)
				movieShowIds := []int{}
				_, ok := theatreWiseAvalibilities[theatreId]
				if ok {
					movieShowIds = []int{}
					a := theatreWiseAvalibilities[theatreId]
					for _, showDetail := range a.MovieShowDetail {
						movieShowIds = append(movieShowIds, showDetail.MovieShowId)
					}
					if ok, _ := util.Contains(movieShowIds, theatreStat.MovieShowId); !ok {
						a.MovieShowDetail = append(a.MovieShowDetail, models.MovieShowDetails{
							MovieShowId: theatreStat.MovieShowId,
							Timing:      slot,
							Seats:       theatreStat.TheatreSeatIds,
						})
						theatreWiseAvalibilities[theatreId] = a
					}
				} else {
					a := models.Availability{
						TheatreId:       theatreStat.TheatreId,
						TheatreName:     theatreStat.TheatreName,
						MovieShowDetail: []models.MovieShowDetails{},
					}
					a.MovieShowDetail = append(a.MovieShowDetail, models.MovieShowDetails{
						MovieShowId: theatreStat.MovieShowId,
						Timing:      slot,
						Seats:       theatreStat.TheatreSeatIds,
					})
					theatreWiseAvalibilities[theatreId] = a
				}

				movieShowIds = []int{}
				if a, ok := theatreWiseAvalibilities[theatreId]; ok {
					for _, movieShowDetail := range a.MovieShowDetail {
						movieShowIds = append(movieShowIds, movieShowDetail.MovieShowId)
					}
				}

				for _, bStat := range bStats {
					if ok, idx := util.Contains(movieShowIds, bStat.MovieShowId); ok {
						if a, ok := theatreWiseAvalibilities[theatreId]; ok {
							movieShowDetail := a.MovieShowDetail[idx]
							seats := []int{}
							seats = append(seats, movieShowDetail.Seats...)
							if yes, seatIdx := util.Contains(seats, bStat.BookedSeatId); yes {
								seats = append(seats[:seatIdx], seats[seatIdx+1:]...)
								movieShowDetail.Seats = []int{}
								movieShowDetail.Seats = append(movieShowDetail.Seats, seats...)
								a.MovieShowDetail[idx] = movieShowDetail
								theatreWiseAvalibilities[theatreId] = a
							}
						}
					}
				}
			}
		} else {
			// no booking on this theatre yet
			for _, theatreStat := range tstats {
				slot := fmt.Sprintf("%d:%d-%d:%d", theatreStat.ShowFromHr, theatreStat.ShowFromMin, theatreStat.ShowTillHr, theatreStat.ShowTillMin)
				a, ok := theatreWiseAvalibilities[theatreId]
				if ok {
					a.MovieShowDetail = append(a.MovieShowDetail, models.MovieShowDetails{
						MovieShowId: theatreStat.MovieShowId,
						Timing:      slot,
						Seats:       theatreStat.TheatreSeatIds,
					})
					theatreWiseAvalibilities[theatreId] = a
				} else {
					a = models.Availability{
						TheatreId:       theatreStat.TheatreId,
						TheatreName:     theatreStat.TheatreName,
						MovieShowDetail: []models.MovieShowDetails{},
					}
					a.MovieShowDetail = append(a.MovieShowDetail, models.MovieShowDetails{
						MovieShowId: theatreStat.MovieShowId,
						Timing:      slot,
						Seats:       theatreStat.TheatreSeatIds,
					})
					theatreWiseAvalibilities[theatreId] = a
				}
			}
		}
	}
	availabilities := []models.Availability{}
	for _, availability := range theatreWiseAvalibilities {
		availabilities = append(availabilities, availability)
	}
	return availabilities, nil
}

type Seats []int

func (s Seats) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

func (s *Seats) Scan(src interface{}) (err error) {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("scan source is not []byte")
	}
	return json.Unmarshal(b, &s)
}

type stat struct {
	BookedSeatId   int    `db:"booked_seat_id"`
	BookingId      int    `db:"booking_id"`
	MovieShowId    int    `db:"movie_show_id"`
	TheatreId      int    `db:"theatre_id"`
	TheatreName    string `db:"theatre_name"`
	ShowId         int    `db:"show_id"`
	ShowFromHr     int    `db:"show_from_hr"`
	ShowFromMin    int    `db:"show_from_min"`
	ShowTillHr     int    `db:"show_till_hr"`
	ShowTillMin    int    `db:"show_till_min"`
	TheatreSeatIds Seats  `db:"theatre_seat_ids"`
}
