package models

type Roles struct {
	Id   int64  `json:"id" db:"Id"`
	Name string `json:"name" db:"Name"`
}
