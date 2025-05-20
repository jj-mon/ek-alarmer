package configmanager

import (
	"bytes"
	"fmt"
	"kuiper-conf/configurator"
	"kuiper-conf/models"
	"log"
	"strings"
	"text/template"
)

type Manager struct {
	cfgr      *configurator.Configurator
	actions   []map[string]any
	projectID string
}

func New(c *configurator.Configurator, projectID string) *Manager {
	return &Manager{
		cfgr: c,
		actions: []map[string]any{
			{
				"sql": map[string]any{
					"url":    "postgres://postgres:postgres@edgex-postgres:5432/edgex_db?sslmode=disable",
					"table":  "core_data.alarm",
					"fields": []string{"project_id", "value", "source_name"},
				},
			},
		},
		projectID: projectID,
	}
}

func (m *Manager) ConfigureDevice(device models.Device) error {
	if err := m.cfgr.DeleteAllRules(); err != nil {
		log.Printf("failed to delete all rules: %s", err)
	}

	if err := m.cfgr.DeleteAllStreams(); err != nil {
		log.Printf("failed to delete all streams: %s", err)
	}

	f, l, _ := strings.Cut(device.Name, ".")
	streamName := fmt.Sprintf("%s%sStream", f, l)
	topicName := fmt.Sprintf("+/+/+/+/+/%s/#", device.Name)

	if err := m.cfgr.CreateStream(streamName, topicName); err != nil {
		log.Printf("failed to create stream: %s", err)
		return err
	}

	for _, source := range device.Sources {
		if err := m.cfgr.CreateRule(m.toRule(source, streamName)); err != nil {
			log.Printf("failed to create rule: %s", err)
			return err
		}
	}

	return nil
}

func (m *Manager) toRule(source models.Source, streamName string) models.Rule {
	t, err := template.New("tmplForRule").Parse(RuleTmpl)
	if err != nil {
		log.Printf("failed to parse tmpl for rule: %s", err)
		return models.Rule{}
	}

	params := models.RuleParams{
		ProjectID:  m.projectID,
		StreamName: streamName,
		SourceName: source.Name,
		LoLo:       source.LoLo,
		Lo:         source.Lo,
		Hi:         source.Hi,
		HiHi:       source.HiHi,
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, params); err != nil {
		log.Printf("failed to execute tmpl for rule: %s", err)
		return models.Rule{}
	}

	return models.Rule{
		ID:      fmt.Sprintf("%sRule", source.Name),
		SQL:     buf.String(),
		Actions: m.actions,
	}
}
