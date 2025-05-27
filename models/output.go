package models

type Threshold struct {
	ID   string `json:"id"`
	HiHi string `json:"hihi"`
	Hi   string `json:"hi"`
	Lo   string `json:"lo"`
	LoLo string `json:"lolo"`
}
