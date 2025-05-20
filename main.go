package main

import (
	"fmt"
	"kuiper-conf/client"
	"kuiper-conf/configmanager"
	"kuiper-conf/configurator"
	"kuiper-conf/models"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	host = "localhost"
	port = 59720
)

type Config struct {
	Devices []models.Device
}

func main() {
	var config models.Config
	if len(os.Args) < 2 {
		config = testConfig()
	} else {
		filePath := os.Args[1]
		config = parseFile(filePath)
	}

	client := client.New(fmt.Sprintf("http://%s:%d", host, port))

	cfgr := configurator.New(client)
	// Создаем плагин sink для записи данных в базу SQL
	if err := cfgr.CreateSinkPlugin("sql", "https://packages.emqx.net/kuiper-plugins/v2.1.3/alpine/sinks/sql_arm64.zip"); err != nil {
		log.Printf("failed to create sink plugin: %s", err)
	}

	projectID := "project1"

	mgr := configmanager.New(cfgr, projectID)

	if err := cfgr.DeleteAllRules(); err != nil {
		log.Printf("failed to delete all rules: %s", err)
	}

	if err := cfgr.DeleteAllStreams(); err != nil {
		log.Printf("failed to delete all streams: %s", err)
	}

	for _, device := range config.Devices {
		if err := mgr.ConfigureDevice(device); err != nil {
			log.Printf("failed to configure device: %s", err)
			return
		}
	}
}

func parseFile(filePath string) models.Config {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("read file error: %v", err)
	}

	var config models.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("unmarshal error: %v", err)
	}

	for i, device := range config.Devices {
		for j := range device.Sources {
			config.Devices[i].Sources[j].HiHi = "17000"
			config.Devices[i].Sources[j].Hi = "14000"
			config.Devices[i].Sources[j].Lo = "10000"
			config.Devices[i].Sources[j].LoLo = "4000"
		}
	}

	return config
}

func testConfig() models.Config {
	return models.Config{
		Devices: []models.Device{
			{
				Name: "modbus.device_505",
				Sources: []models.Source{
					{
						Name: "R_1",
						HiHi: "17000",
						Hi:   "14000",
						Lo:   "10000",
						LoLo: "4000",
					},
					{
						Name: "R_2",
						HiHi: "17000",
						Hi:   "14000",
						Lo:   "10000",
						LoLo: "4000",
					},
					{
						Name: "R_3",
						HiHi: "17000",
						Hi:   "14000",
						Lo:   "10000",
						LoLo: "4000",
					},
					{
						Name: "R_4",
						HiHi: "17000",
						Hi:   "14000",
						Lo:   "10000",
						LoLo: "4000",
					},
					{
						Name: "R_5",
						HiHi: "17000",
						Hi:   "14000",
						Lo:   "10000",
						LoLo: "4000",
					},
				},
			},
		},
	}
}
