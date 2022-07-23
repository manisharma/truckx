package models

type Seats struct {
	Id        int64 `json:"-" db:"Id"`
	TheatreId int64 `json:"-" db:"TheatreId"`
}
