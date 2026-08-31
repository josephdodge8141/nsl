package nsl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const apiPath = "/api/v2"

type Client struct {
	BaseURL  string
	HTTP     *http.Client
	APIToken string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: strings.TrimRight(baseURL, "/"), HTTP: &http.Client{Timeout: 30 * time.Second}, APIToken: os.Getenv("NSL_API_TOKEN")}
}

func (c *Client) FetchVersion() (string, error) {
	var response struct {
		Version string `json:"version"`
	}
	if _, err := c.do(http.MethodGet, apiPath+"/version", nil, "", &response); err != nil {
		return "", err
	}
	return response.Version, nil
}

func (c *Client) LocalNode() (Node, error) {
	var node Node
	_, err := c.do(http.MethodGet, apiPath+"/node", nil, "", &node)
	return node, err
}

func (c *Client) Config() (Config, error) {
	var registryConfig Config
	_, err := c.do(http.MethodGet, "/api/config", nil, "", &registryConfig)
	return registryConfig, err
}

func (c *Client) Nodes() ([]Node, error) {
	var nodes []Node
	_, err := c.do(http.MethodGet, apiPath+"/nodes", nil, "", &nodes)
	return nodes, err
}

func (c *Client) List() ([]App, error) {
	var apps []App
	_, err := c.do(http.MethodGet, apiPath+"/apps", nil, "", &apps)
	return apps, err
}

func (c *Client) Create(input AppInput) (App, error) {
	var app App
	_, err := c.do(http.MethodPost, apiPath+"/apps", input, "", &app)
	return app, err
}

func (c *Client) Delete(app App) error {
	_, err := c.do(http.MethodDelete, apiPath+"/apps/"+app.ID, nil, appETag(app), nil)
	return err
}

func (c *Client) do(method, path string, input any, ifMatch string, output any) (http.Header, error) {
	var body io.Reader
	if input != nil {
		encoded, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(encoded)
	}
	request, err := http.NewRequest(method, c.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	if input != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if ifMatch != "" {
		request.Header.Set("If-Match", ifMatch)
	}
	if c.APIToken != "" {
		request.Header.Set("Authorization", "Bearer "+c.APIToken)
	}
	response, err := c.HTTP.Do(request)
	if err != nil {
		return nil, fmt.Errorf("registry unavailable: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		message, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		return response.Header, fmt.Errorf("registry returned %d: %s", response.StatusCode, strings.TrimSpace(string(message)))
	}
	if output != nil {
		if err := json.NewDecoder(response.Body).Decode(output); err != nil {
			return response.Header, fmt.Errorf("decode registry response: %w", err)
		}
	}
	return response.Header, nil
}

func appETag(app App) string {
	return fmt.Sprintf(`"app:%s:%d"`, app.ID, app.Generation)
}
