package models

import "time"

type Booking struct {
	Id            int64     `json:"id" db:"id"`
	MoviesShowsId int64     `json:"moviesShowsId" db:"movies_shows_id"`
	CreationDate  time.Time `json:"creationDate" db:"creation_date"`
	BookingDate   time.Time `json:"bookingDate" db:"booking_date"`
}

type CreateBooking struct {
	TheatreId   int           `json:"theatreId"`
	State       string        `json:"state"`
	MovieId     int           `json:"movieId"`
	MovieShowId int           `json:"movieShowId"`
	Seats       []BookingSeat `json:"seats"`
	DateStr     string        `json:"date"`
	Date        time.Time     `json:"-"`
	UserId      int           `json:"-"`
}
