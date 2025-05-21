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
					"fields": []string{"project_id", "value", "source_name", "alarm_level"},
				},
			},
		},
		projectID: projectID,
	}
}

func (m *Manager) ConfigureDevice(device models.Device) (models.Stream, error) {
	var stream models.Stream

	f, l, _ := strings.Cut(device.Name, ".")
	streamName := fmt.Sprintf("%s%sStream", f, l)
	topicName := fmt.Sprintf("+/+/+/+/+/%s/#", device.Name)

	stream.Name = streamName
	stream.SQL = fmt.Sprintf("CREATE stream %s () WITH (FORMAT=\"JSON\", DATASOURCE=\"%s\", SHARED=\"true\")", streamName, topicName)

	log.Printf("stream %s configured", streamName)

	for _, source := range device.Sources {
		stream.Rules = append(stream.Rules, m.toRule(source, streamName))
		log.Printf("rule %s configured", source.Name)
	}

	return stream, nil
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
