package main

import (
	"fmt"
	"kuiper-conf/client"
	"kuiper-conf/configmanager"
	"kuiper-conf/configurator"
	"kuiper-conf/models"
	"log"
)

const (
	host = "localhost"
	port = 59720
)

func main() {
	client := client.New(fmt.Sprintf("http://%s:%d", host, port))

	cfgr := configurator.New(client)
	// Создаем плагин sink для записи данных в базу SQL
	if err := cfgr.CreateSinkPlugin("sql", "https://packages.emqx.net/kuiper-plugins/v2.1.3/alpine/sinks/sql_arm64.zip"); err != nil {
		log.Printf("failed to create sink plugin: %s", err)
	}

	projectID := "project1"

	mgr := configmanager.New(cfgr, projectID)

	device := models.Device{
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
	}

	if err := mgr.ConfigureDevice(device); err != nil {
		log.Printf("failed to configure device: %s", err)
		return
	}
}
