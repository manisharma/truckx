package models

type Theatre struct {
	Id         int64  `json:"id" db:"Id"`
	Name       string `json:"name" db:"Name"`
	Address    string `json:"address" db:"Address"`
	State      string `json:"state" db:"State"`
	PinCode    int    `json:"pinCode" db:"PinCode"`
	TotalSeats int    `json:"totalSeats" db:"TotalSeats"`
}
