package models

type Shows struct {
	Id        int `json:"id" db:"Id"`
	TheatreId int `json:"theatreId" db:"TheatreId"`
	FromHr    int `json:"fromHr" db:"FromHr"`
	FromMin   int `json:"fromMin" db:"FromMin"`
	TillHr    int `json:"tillHr" db:"TillHr"`
	TillMin   int `json:"tillMin" db:"TillMin"`
}
