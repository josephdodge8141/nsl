package nsl

import "time"

type App struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	AppType          string    `json:"app_type"`
	RouteRule        string    `json:"route_rule"`
	TargetURL        string    `json:"target_url,omitempty"`
	DocsURL          string    `json:"docs_url,omitempty"`
	ConnectionString string    `json:"connection_string,omitempty"`
	ContainerName    string    `json:"container_name"`
	NoAuth           bool      `json:"no_auth"`
	Enabled          bool      `json:"enabled"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
