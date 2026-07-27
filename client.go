package nsl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTP:    http.DefaultClient,
	}
}

const apiPath = "/api/v1"

func (c *Client) FetchVersion() (string, error) {
	resp, err := c.HTTP.Get(c.BaseURL + apiPath + "/version")
	if err != nil {
		return "", fmt.Errorf("fetch version: %w", err)
	}
	defer resp.Body.Close()

	var v struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return "", fmt.Errorf("fetch version: decode: %w", err)
	}
	return v.Version, nil
}

func (c *Client) List() ([]App, error) {
	resp, err := c.HTTP.Get(c.BaseURL + apiPath + "/apps")
	if err != nil {
		return nil, fmt.Errorf("list apps: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list apps: %s", resp.Status)
	}

	var apps []App
	if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
		return nil, fmt.Errorf("list apps: decode: %w", err)
	}
	return apps, nil
}

func (c *Client) Create(app App) (*App, error) {
	body, err := json.Marshal(app)
	if err != nil {
		return nil, fmt.Errorf("create app: marshal: %w", err)
	}

	resp, err := c.HTTP.Post(c.BaseURL+apiPath+"/apps", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create app: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create app: %s: %s", resp.Status, strings.TrimSpace(string(raw)))
	}

	var created App
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, fmt.Errorf("create app: decode: %w", err)
	}
	return &created, nil
}

func (c *Client) Delete(id string) error {
	req, err := http.NewRequest(http.MethodDelete, c.BaseURL+apiPath+"/apps/"+id, nil)
	if err != nil {
		return fmt.Errorf("delete app: %w", err)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("delete app: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete app: %s", resp.Status)
	}
	return nil
}
