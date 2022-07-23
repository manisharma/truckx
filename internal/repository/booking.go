package repository

import (
	"context"
	"fmt"
	"strings"
	"truckx/internal/models"

	"github.com/jmoiron/sqlx"
)

type BookingRepo interface {
	Create(ctx context.Context, createBooking models.CreateBooking) (int, error)
	Get(ctx context.Context, bookingId int) ([]models.BookingSeat, error)
	List(ctx context.Context, userId int) ([]models.BookingSeat, error)
}

func NewBookingRepo(db *sqlx.DB) BookingRepo {
	return &bookingImpl{db}
}

type bookingImpl struct {
	db *sqlx.DB
}

func (r bookingImpl) Create(ctx context.Context, createBooking models.CreateBooking) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, nil
	}
	success := false
	defer func() {
		if success {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	// insert booking
	row, err := tx.QueryContext(ctx,
		"INSERT INTO bookings(movies_shows_id, creation_date, booking_date) VALUES ($1, CURRENT_TIMESTAMP, $2); RETURNING id",
		createBooking.MovieShowId, createBooking.Date)
	if err != nil {
		return 0, err
	}
	var bookingId int
	err = row.Scan(&bookingId)
	if err != nil {
		return 0, err
	}

	// insert user_bookings, only if user is logged in
	if createBooking.UserId != 0 {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO user_bookings(user_id, booking_id) VALUES ($1, $2)",
			createBooking.UserId, bookingId)
		if err != nil {
			return 0, err
		}
	}

	// bulk insert booking details
	valueStrings := make([]string, 0, len(createBooking.Seats))
	valueArgs := make([]interface{}, 0, len(createBooking.Seats)*5)
	i := 0
	for _, seat := range createBooking.Seats {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
		valueArgs = append(valueArgs, bookingId)
		valueArgs = append(valueArgs, seat.SeatId)
		valueArgs = append(valueArgs, seat.GuestName)
		valueArgs = append(valueArgs, seat.GuestEmail)
		valueArgs = append(valueArgs, seat.GuestPhone)
	}
	_, err = tx.ExecContext(ctx,
		fmt.Sprintf("INSERT INTO bookings_seats(booking_id, seat_id, guest_name, guest_email, guest_phone) VALUES %s", strings.Join(valueStrings, ",")), valueArgs...)
	if err != nil {
		return 0, err
	}
	success = true
	return bookingId, nil
}

func (r bookingImpl) Get(ctx context.Context, bookingId int) ([]models.BookingSeat, error) {
	details := []models.BookingSeat{}
	rows, err := r.db.QueryxContext(ctx, "SELECT * FROM bookings_seats WHERE booking_id = $1", bookingId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	detail := models.BookingSeat{}
	for rows.Next() {
		err := rows.StructScan(&detail)
		if err != nil {
			return nil, err
		}
		details = append(details, detail)
	}
	return details, nil
}

func (r bookingImpl) List(ctx context.Context, userId int) ([]models.BookingSeat, error) {
	details := []models.BookingSeat{}
	rows, err := r.db.QueryxContext(ctx, "SELECT bs.* FROM bookings_seats bs JOIN user_bookings ub ON ub.booking_id = bs.booking_id WHERE ub.user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	detail := models.BookingSeat{}
	for rows.Next() {
		err := rows.StructScan(&detail)
		if err != nil {
			return nil, err
		}
		details = append(details, detail)
	}
	return details, nil
}
