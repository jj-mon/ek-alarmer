package main

import (
	"encoding/json"
	"fmt"
	"kuiper-conf/client"
	"kuiper-conf/configmanager"
	"kuiper-conf/configurator"
	"kuiper-conf/devparser"
	"kuiper-conf/models"
	"log"
	"os"
	"runtime"
)

const (
	host = "localhost"
	port = 59720
)

type Config struct {
	Devices []models.Device
}

func main() {
	if len(os.Args) < 2 {
		log.Printf("input path to device file")
		return
	}

	filePath := os.Args[1]
	// парсим файл edgex девайса
	tresholds := devparser.Parse(filePath)
	// записываем граничные значения для каждого девайса и для каждого сурса в виде JSON
	writeTresholdsJSON(tresholds)
	// создаем HTTP клиента для взаимодействия с REST API eKuiper
	client := client.New(fmt.Sprintf("http://%s:%d", host, port))
	// создаем конфигуратор eKupier
	cfgr := configurator.New(client)
	// Создаем плагин sink для записи данных в базу SQL
	arch := runtime.GOARCH
	pluginURL := fmt.Sprintf("https://packages.emqx.net/kuiper-plugins/v2.1.3/alpine/sinks/sql_%s.zip", arch)
	// Создаем SQL плагин для слива данных в БД в actions eKuiper's
	if err := cfgr.CreateSinkPlugin("sql", pluginURL); err != nil {
		log.Printf("failed to create sink plugin: %s", err)
	}

	projectID := "project1"
	// Создаем конфигуратор менеджера для eKuiper'а
	mgr := configmanager.New(cfgr, projectID)
	// На всякий чистим все правила
	if err := cfgr.DeleteAllRules(); err != nil {
		log.Printf("failed to delete all rules: %s", err)
	}
	// Чистим все потоки
	if err := cfgr.DeleteAllStreams(); err != nil {
		log.Printf("failed to delete all streams: %s", err)
	}
	// С помощью менеджера конфига екупера создаем таблицу для поиска граничных значений для сурса
	mgr.CreateTresholdsTable("./tresholds.json")

	

	_ = mgr
}

func writeTresholdsJSON(tresholds []models.Threshold) {
	file, err := os.Create("tresholds.json")
	if err != nil {
		log.Fatalf("failed create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(tresholds); err != nil {
		log.Fatalf("failed encode tresholds: %v", err)
	}

	log.Print("file success created: tresholds.json")
}
