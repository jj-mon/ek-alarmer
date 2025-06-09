package main

import (
	"encoding/json"
	"fmt"
	"kuiper-conf/client"
	"kuiper-conf/configurator"
	"kuiper-conf/models"
	"log"
	"os"
	"runtime"
)

const (
	host = "localhost"
	port = 59720
)

var ruleset = `
{
   "streams":{
      "eventStream":"CREATE STREAM eventStream () WITH (DATASOURCE=\"event\", FORMAT=\"json\", TYPE=\"memory\", SHARED=\"false\");",
      "mainStream":"CREATE STREAM mainStream () WITH (DATASOURCE=\"edgex/events/#\", FORMAT=\"json\", TYPE=\"mqtt\", SHARED=\"false\");"
   },
   "tables":{
      "threshold":"CREATE TABLE threshold (project_id string, source_name string, lolo bigint, lo bigint, hi bigint, hihi bigint) WITH (DATASOURCE=\"core_data.threshold_table\", FORMAT=\"json\", CONF_KEY=\"postgresql_config\", TYPE=\"sql\", KIND=\"lookup\", KEY=\"source_name\");"
   },
   "rules":{
      "rule_1":"{\"triggered\":true,\"id\":\"rule_1\",\"name\":\"take data from main stream and send it to memory topic\",\"sql\":\"SELECT \\n    mainStream.payload.event.readings[0].deviceName as deviceName,\\n    mainStream.payload.event.readings[0].profileName as profileName,\\n    mainStream.payload.event.readings[0].resourceName as resourceName,\\n    mainStream.payload.event.readings[0].value as cur_val\\nFROM mainStream \",\"actions\":[{\"memory\":{\"bufferLength\":1024,\"enableCache\":false,\"format\":\"json\",\"omitIfEmpty\":false,\"runAsync\":false,\"sendSingle\":true,\"topic\":\"event\"}}],\"options\":{\"debug\":false,\"isEventTime\":false,\"lateTolerance\":\"1s\",\"concurrency\":1,\"bufferLength\":1024,\"sendMetaToSink\":false,\"sendNilField\":false,\"sendError\":false,\"checkpointInterval\":\"5m0s\",\"restartStrategy\":{\"delay\":\"1s\",\"multiplier\":2,\"maxDelay\":\"30s\",\"jitterFactor\":0.1}}}",
      "rule_2":"{\"triggered\":true,\"id\":\"rule_2\",\"name\":\"take event from eventStream and work with it\",\"sql\":\"SELECT \\n    threshold.project_id as project_id,\\n    eventStream.resourceName as source_name,\\n    eventStream.cur_val as cur_val,\\n    lag(eventStream.cur_val) OVER (PARTITION BY eventStream.resourceName) as prev_val,\\n    CASE\\n        WHEN lag(cast(eventStream.cur_val, 'bigint')) OVER (PARTITION BY eventStream.resourceName) \\u003e threshold.lolo\\n            AND cast(eventStream.cur_val, 'bigint') \\u003c= threshold.lolo THEN \\\"CRITICAL_LOW\\\"\\n        WHEN lag(cast(eventStream.cur_val, 'bigint')) OVER (PARTITION BY eventStream.resourceName) \\u003e threshold.lo\\n            AND cast(eventStream.cur_val, 'bigint') BETWEEN threshold.lo AND threshold.lolo THEN \\\"WARNING_LOW\\\"\\n        WHEN lag(cast(eventStream.cur_val, 'bigint')) OVER (PARTITION BY eventStream.resourceName) \\u003c threshold.hi\\n            AND cast(eventStream.cur_val, 'bigint') BETWEEN threshold.hi AND threshold.hihi THEN \\\"WARNING_HI\\\"\\n        WHEN lag(cast(eventStream.cur_val, 'bigint')) OVER (PARTITION BY eventStream.resourceName) \\u003c threshold.hihi\\n            AND cast(eventStream.cur_val, 'bigint') \\u003e= threshold.hihi THEN \\\"CRITICAL_HI\\\"\\n    END AS status\\nFROM eventStream\\nINNER JOIN threshold \\nON eventStream.resourceName = threshold.source_name\\nWHERE isNull(status) = false\",\"actions\":[{\"mqtt\":{\"bufferLength\":1024,\"enableCache\":false,\"format\":\"json\",\"omitIfEmpty\":false,\"runAsync\":false,\"sendSingle\":true,\"server\":\"tcp://edgex-mqtt-broker:1883\",\"topic\":\"result\"}},{\"sql\":{\"bufferLength\":1024,\"dburl\":\"postgres://postgres:postgres@edgex-postgres:5432/edgex_db?sslmode=disable\",\"enableCache\":false,\"fields\":[\"project_id\",\"source_name\",\"cur_val\",\"prev_val\",\"status\"],\"format\":\"json\",\"omitIfEmpty\":false,\"runAsync\":false,\"sendSingle\":true,\"table\":\"core_data.alert\"}}],\"options\":{\"debug\":false,\"isEventTime\":false,\"lateTolerance\":\"1s\",\"concurrency\":1,\"bufferLength\":1024,\"sendMetaToSink\":false,\"sendNilField\":false,\"sendError\":false,\"checkpointInterval\":\"5m0s\",\"restartStrategy\":{\"delay\":\"1s\",\"multiplier\":2,\"maxDelay\":\"30s\",\"jitterFactor\":0.1}}}"
   }
}
`

var sourceConfig = map[string]any{
	"dburl": "postgres://postgres:postgres@edgex-postgres:5432/edgex_db?sslmode=disable",
	"internalSqlQueryCfg": map[string]string{
		"table": "core_data.threshold_table",
	},
}

func main() {
	// записываем граничные значения для каждого девайса и для каждого сурса в виде JSON
	// создаем HTTP клиента для взаимодействия с REST API eKuiper
	client := client.New(fmt.Sprintf("http://%s:%d", host, port))
	// создаем конфигуратор eKupier
	cfgr := configurator.New(client)
	// Создаем плагин sink для записи данных в базу SQL
	arch := runtime.GOARCH
	pluginURL := fmt.Sprintf("https://packages.emqx.net/kuiper-plugins/v2.1.3/alpine/sinks/sql_%s.zip", arch)
	// Создаем SQL плагин для слива данных в БД в actions eKuiper's
	if err := cfgr.CreateSinkPlugin("sql", pluginURL); err != nil {
		log.Printf("failed to create sink plugin: %v", err)
	}
	pluginURL = fmt.Sprintf("https://packages.emqx.net/kuiper-plugins/v2.1.3/alpine/sources/sql_%s.zip", arch)
	if err := cfgr.CreateSourcePlugin("sql", pluginURL); err != nil {
		log.Printf("failed to create source plugin: %v", err)
	}

	if err := cfgr.RegisterSrcConfig(sourceConfig); err != nil {
		log.Printf("failed to register source config: %v", err)
	}

	if err := cfgr.CreateRuleset(ruleset); err != nil {
		log.Printf("failed to create ruleset: %v", err)
	}
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
