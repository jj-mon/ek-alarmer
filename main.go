package main

import (
	"encoding/json"
	"fmt"
	"kuiper-conf/client"
	"kuiper-conf/configmanager"
	"kuiper-conf/configurator"
	"kuiper-conf/models"
	"log"
	"os"
	"runtime"

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
	arch := runtime.GOARCH
	pluginURL := fmt.Sprintf("https://packages.emqx.net/kuiper-plugins/v2.1.3/alpine/sinks/sql_%s.zip", arch)
	if err := cfgr.CreateSinkPlugin("sql", pluginURL); err != nil {
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

	// инициализируем рулсет перед его заполнением
	ruleset := models.Ruleset{
		Streams: map[string]string{},
		Rules:   map[string]string{},
		Tables:  map[string]string{},
	}
	for _, device := range config.Devices {
		stream, err := mgr.ConfigureDevice(device)
		if err != nil {
			log.Printf("failed to configure device: %s", err)
			return
		}
		ruleset.Streams[stream.Name] = stream.SQL
		for _, rule := range stream.Rules {
			b, err := json.Marshal(rule)
			if err != nil {
				log.Printf("failed to configure device: %s", err)
				continue
			}

			ruleset.Rules[rule.ID] = string(b)
		}
	}

	// оправляем рулсет в екупер
	if err := cfgr.CreateRuleset(ruleset); err != nil {
		log.Printf("failed to create ruleset: %s", err)
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
