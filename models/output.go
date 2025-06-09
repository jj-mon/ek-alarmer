package models

type Threshold struct {
	ProjectID  string `json:"project_id"`
	SourceName string `json:"source_name"`
	HiHi       string `json:"hihi"`
	Hi         string `json:"hi"`
	Lo         string `json:"lo"`
	LoLo       string `json:"lolo"`
}
