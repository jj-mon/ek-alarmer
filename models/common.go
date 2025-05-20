package models

type Rule struct {
	ID      string           `json:"id"`
	SQL     string           `json:"sql"`
	Actions []map[string]any `json:"actions"`
}

type RuleParams struct {
	ProjectID  string
	StreamName string
	SourceName string
	LoLo       string
	Lo         string
	Hi         string
	HiHi       string
}
