package configurator

import (
	"encoding/json"
	"kuiper-conf/client"
	"log"
)

type Configurator struct {
	client *client.Client
}

func New(c *client.Client) *Configurator {
	return &Configurator{client: c}
}

func (c *Configurator) CreateSinkPlugin(name, url string) error {
	data := map[string]string{
		"name": name,
		"file": url,
	}

	_, err := c.client.DoPOST("/plugins/sinks", data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Configurator) CreateSourcePlugin(name, url string) error {
	data := map[string]string{
		"name": name,
		"file": url,
	}

	_, err := c.client.DoPOST("/plugins/sources", data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Configurator) RegisterSrcConfig(config map[string]any) error {
	resp, err := c.client.DoPUT("/metadata/sources/sql/confKeys/postgresql_config", config)
	if err != nil {
		log.Printf("error register source config: %s", resp)
		return err
	}

	return nil
}

func (c *Configurator) CreateRuleset(ruleset string) error {
	data := map[string]string{
		"content": string(ruleset),
	}

	resp, err := c.client.DoPOST("/ruleset/import", data)
	if err != nil {
		log.Printf("error creating ruleset: %s", resp)
		return err
	}

	return nil
}

func (c *Configurator) DeleteAllStreams() error {
	var streams []string

	resp, err := c.client.DoGET("/streams")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp, &streams); err != nil {
		return err
	}

	for _, stream := range streams {
		if err := c.DropStream(stream); err != nil {
			return err
		}
		log.Printf("stream %s dropped", stream)
	}

	return nil
}

func (c *Configurator) DropStream(id string) error {
	_, err := c.client.DoDELETE("/streams", id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Configurator) DeleteAllRules() error {
	var rules []map[string]any

	resp, err := c.client.DoGET("/rules")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp, &rules); err != nil {
		return err
	}

	for _, rule := range rules {
		if err := c.DropRule(rule["id"].(string)); err != nil {
			return err
		}
		log.Printf("rule %s dropped", rule["id"].(string))
	}

	return nil
}

func (c *Configurator) DropRule(id string) error {
	_, err := c.client.DoDELETE("/rules", id)
	if err != nil {
		return err
	}

	return nil
}
