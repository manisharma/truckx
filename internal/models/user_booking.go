package models

type UserBooking struct {
	UserId    int64 `json:"userId" db:"UserId"`
	BookingId int64 `json:"bookingId" db:"BookingId"`
}
