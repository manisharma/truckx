package models

type Availability struct {
	TheatreId       int                `json:"theatreId"`
	TheatreName     string             `json:"theatreName"`
	MovieShowDetail []MovieShowDetails `json:"movieShowDetail"`
}

type MovieShowDetails struct {
	MovieShowId int    `json:"movieShowId"`
	Timing      string `json:"timing"`
	Seats       []int  `json:"seats"`
}

// type Availability struct {
// 	TheatreId   int    `json:"theatreId"`
// 	TheatreName string `json:"theatreName"`
// 	Slot        string `json:"slot"`
// 	Seats       []int  `json:"seats"`
// }
