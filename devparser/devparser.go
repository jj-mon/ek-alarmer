package devparser

import (
	"fmt"
	"kuiper-conf/models"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func Parse(filePath string) []models.Threshold {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("read file error: %v", err)
	}

	var config models.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("unmarshal error: %v", err)
	}

	return createTresholds(config)
}

func createTresholds(config models.Config) []models.Threshold {
	var tresholds []models.Threshold

	for _, device := range config.Devices {
		for _, source := range device.Sources {
			treshold := models.Threshold{
				ProjectID:  "ID-12345",
				SourceName: fmt.Sprintf("%s_%s", device.Name, source.Name),
				HiHi:       "17000",
				Hi:         "14000",
				Lo:         "10000",
				LoLo:       "4000",
			}
			tresholds = append(tresholds, treshold)
		}
	}

	return tresholds
}
