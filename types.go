// Package nsl implements the client for the distributed NSL registry.
package nsl

import "time"

// AuthPolicy controls whether NSL or the target service authenticates a route.
type AuthPolicy string

const (
	// AuthBrowser protects browser routes with Keycloak.
	AuthBrowser AuthPolicy = "browser"
	// AuthUpstream leaves authentication to the target service.
	AuthUpstream AuthPolicy = "upstream"
)

// Node identifies an enrolled Not-So-Localhost machine.
type Node struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Slug  string   `json:"slug"`
	Roles []string `json:"roles,omitempty"`
}

// Config describes the local registry and node identity.
type Config struct {
	Domain   string `json:"domain"`
	NodeID   string `json:"node_id"`
	NodeName string `json:"node_name"`
	NodeSlug string `json:"node_slug"`
}

// Route describes one public routing rule for an application.
type Route struct {
	ID       string     `json:"id"`
	Rule     string     `json:"rule"`
	Priority int        `json:"priority"`
	Auth     AuthPolicy `json:"auth"`
}

// RouteInput contains mutable route fields sent to the registry.
type RouteInput struct {
	Rule     string     `json:"rule"`
	Priority int        `json:"priority,omitempty"`
	Auth     AuthPolicy `json:"auth"`
}

// App is an HTTP service registered to one node.
type App struct {
	ID          string    `json:"id"`
	NodeID      string    `json:"node_id"`
	Generation  uint64    `json:"generation"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	TargetURL   string    `json:"target_url"`
	PublicURL   string    `json:"public_url"`
	Routes      []Route   `json:"routes"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AppInput contains mutable application fields sent to the registry.
type AppInput struct {
	NodeID      string       `json:"node_id,omitempty"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	TargetURL   string       `json:"target_url"`
	Routes      []RouteInput `json:"routes,omitempty"`
	Enabled     *bool        `json:"enabled,omitempty"`
}
