// Package main tests the nsl command-line client.
package main

import (
	"strings"
	"testing"

	"github.com/josephdodge8141/nsl"
)

func TestLiteLLMPolicyUsesMutuallyExclusiveRoutes(t *testing.T) {
	t.Parallel()
	routes := policyRoutes("litellm", nsl.Node{Slug: "laptop3"}, "joedodge.dev", "litellm")
	if len(routes) != 2 {
		t.Fatalf("got %d routes", len(routes))
	}
	if routes[0].Auth != nsl.AuthBrowser || routes[1].Auth != nsl.AuthUpstream {
		t.Fatalf("unexpected policies: %#v", routes)
	}
	if !strings.Contains(routes[0].Rule, "!(Path(`/v1`)") {
		t.Fatalf("UI route is not excluding /v1: %s", routes[0].Rule)
	}
	if !strings.Contains(routes[1].Rule, "Path(`/v1`)") {
		t.Fatalf("API route does not include /v1: %s", routes[1].Rule)
	}
}

func TestAddRequiresTarget(t *testing.T) {
	t.Parallel()
	err := validateAddFlags(addFlags{Name: "api", Policy: "browser"})
	if err == nil || !strings.Contains(err.Error(), "target-url") {
		t.Fatalf("error = %v", err)
	}
}
