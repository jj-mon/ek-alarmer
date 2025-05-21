package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Client struct {
	url    string
	client *http.Client
}

func New(url string) *Client {
	return &Client{
		url:    url,
		client: &http.Client{},
	}
}

func (c *Client) DoPOST(path string, data any) (string, error) {
	req, err := http.NewRequest("POST", c.url+path, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	log.Printf("DATA SIZE: %d", len(b))

	req.Body = io.NopCloser(bytes.NewBuffer(b))

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return string(respBody), fmt.Errorf("failed do POST request: %s", resp.Status)
	}

	return string(respBody), nil
}

func (c *Client) DoDELETE(path string, id string) (string, error) {
	req, err := http.NewRequest("DELETE", c.url+path+"/"+id, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return string(respBody), fmt.Errorf("failed do POST request: %s", resp.Status)
	}

	return string(respBody), nil
}

func (c *Client) DoGET(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.url+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed do GET request: %s", resp.Status)
	}

	return respBody, nil
}
