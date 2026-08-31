// Package nsl implements the client for the distributed NSL registry.
package nsl

import "time"

type AuthPolicy string

const (
	AuthBrowser  AuthPolicy = "browser"
	AuthUpstream AuthPolicy = "upstream"
)

type Node struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Slug  string   `json:"slug"`
	Roles []string `json:"roles,omitempty"`
}

type Config struct {
	Domain   string `json:"domain"`
	NodeID   string `json:"node_id"`
	NodeName string `json:"node_name"`
	NodeSlug string `json:"node_slug"`
}

type Route struct {
	ID       string     `json:"id"`
	Rule     string     `json:"rule"`
	Priority int        `json:"priority"`
	Auth     AuthPolicy `json:"auth"`
}

type RouteInput struct {
	Rule     string     `json:"rule"`
	Priority int        `json:"priority,omitempty"`
	Auth     AuthPolicy `json:"auth"`
}

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

type AppInput struct {
	NodeID      string       `json:"node_id,omitempty"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	TargetURL   string       `json:"target_url"`
	Routes      []RouteInput `json:"routes,omitempty"`
	Enabled     *bool        `json:"enabled,omitempty"`
}
