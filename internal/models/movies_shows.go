package models

import "time"

type MoviesShows struct {
	Id       int64     `json:"id" db:"Id"`
	MovieId  int64     `json:"movieId" db:"MovieId"`
	ShowId   int64     `json:"showId" db:"ShowId"`
	FromDate time.Time `json:"fromDate" db:"FromDate"`
	TillDate time.Time `json:"tillDate" db:"TillDate"`
}
