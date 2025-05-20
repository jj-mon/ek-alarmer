package models

type Config struct {
	Devices []Device `yaml:"deviceList"`
}

type Device struct {
	Name    string   `yaml:"name"`
	Sources []Source `yaml:"autoEvents"`
}

type Source struct {
	Name string `yaml:"sourceName"`
	HiHi string
	Hi   string
	Lo   string
	LoLo string
}
