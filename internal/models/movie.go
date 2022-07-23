package models

import "time"

type Movie struct {
	Id          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	ReleaseDate time.Time `json:"releaseDate" db:"release_date"`
}
