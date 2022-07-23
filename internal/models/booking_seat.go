package models

type BookingSeat struct {
	BookingId  int64  `json:"bookingId" db:"booking_id"`
	SeatId     int64  `json:"seatId" db:"seat_id"`
	GuestName  string `json:"guestName" db:"guest_name"`
	GuestEmail string `json:"guestEmail" db:"guest_email"`
	GuestPhone string `json:"guestPhone" db:"guest_phone"`
}
