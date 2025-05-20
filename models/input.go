package models

type Device struct {
	Name    string
	Sources []Source
}

type Source struct {
	Name string
	HiHi string
	Hi   string
	Lo   string
	LoLo string
}
