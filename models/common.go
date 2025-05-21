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

type Ruleset struct {
	Streams map[string]string `json:"streams"`
	Rules   map[string]string `json:"rules"`
	Tables  map[string]string `json:"tables"`
}

type Stream struct {
	Name  string
	SQL   string
	Rules []Rule
}
