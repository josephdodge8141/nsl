// Package main implements the nsl command-line client.
package main

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/josephdodge8141/nsl"
)

type addFlags struct {
	Name        string
	TargetURL   string
	Description string
	Policy      string
	Disabled    bool
}

var hostnameRE = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?$`)

func addCmd(apiURL string, flags addFlags) error {
	if err := validateAddFlags(flags); err != nil {
		return usageError(err.Error(), `nsl add --name <name> --target-url <url> [--policy browser|upstream|litellm]`)
	}
	client := nsl.NewClient(apiURL)
	registryConfig, err := client.Config()
	if err != nil {
		return err
	}
	node, err := client.LocalNode()
	if err != nil {
		return err
	}
	nodeID := node.ID
	enabled := !flags.Disabled
	routes := policyRoutes(flags.Name, node, registryConfig.Domain, flags.Policy)
	existingApps, err := client.List()
	if err != nil {
		return err
	}
	for _, existing := range existingApps {
		if existing.NodeID != nodeID || !strings.EqualFold(existing.Name, flags.Name) {
			continue
		}
		desiredRoutes := make([]routeComparable, len(routes))
		for index, route := range routes {
			desiredRoutes[index] = routeComparable{Rule: route.Rule, Priority: route.Priority, Auth: route.Auth}
		}
		existingRoutes := make([]routeComparable, len(existing.Routes))
		for index, route := range existing.Routes {
			existingRoutes[index] = routeComparable{Rule: route.Rule, Priority: route.Priority, Auth: route.Auth}
		}
		sort.Slice(desiredRoutes, func(i, j int) bool { return desiredRoutes[i].Rule < desiredRoutes[j].Rule })
		sort.Slice(existingRoutes, func(i, j int) bool { return existingRoutes[i].Rule < existingRoutes[j].Rule })
		if existing.TargetURL != flags.TargetURL || existing.Description != flags.Description || existing.Enabled != !flags.Disabled || !reflect.DeepEqual(existingRoutes, desiredRoutes) {
			return usageError("app is already registered with different configuration", "Use the apps portal to edit the existing registration")
		}
		fmt.Println("app:")
		fmt.Printf("  id: %s\n", quote(existing.ID))
		fmt.Printf("  name: %s\n", quote(existing.Name))
		fmt.Printf("  node_id: %s\n", quote(existing.NodeID))
		fmt.Printf("  url: %s\n", quote(existing.PublicURL))
		fmt.Println("  status: already_registered")
		return nil
	}
	created, err := client.Create(nsl.AppInput{
		NodeID: nodeID, Name: flags.Name, TargetURL: flags.TargetURL,
		Description: flags.Description, Routes: routes, Enabled: &enabled,
	})
	if err != nil {
		return err
	}
	fmt.Println("app:")
	fmt.Printf("  id: %s\n", quote(created.ID))
	fmt.Printf("  name: %s\n", quote(created.Name))
	fmt.Printf("  node_id: %s\n", quote(created.NodeID))
	fmt.Printf("  url: %s\n", quote(created.PublicURL))
	fmt.Printf("  status: registered\n")
	return nil
}

func validateAddFlags(flags addFlags) error {
	if flags.Name == "" {
		return errors.New("--name is required")
	}
	if !hostnameRE.MatchString(flags.Name) {
		return errors.New("--name must contain only letters, numbers, and internal hyphens")
	}
	if len(flags.Name) > 30 {
		return errors.New("--name must be 30 characters or fewer")
	}
	if flags.TargetURL == "" {
		return errors.New("--target-url is required")
	}
	parsed, err := url.Parse(flags.TargetURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" || parsed.User != nil {
		return errors.New("--target-url must be an HTTP URL without credentials")
	}
	if flags.Policy != "browser" && flags.Policy != "upstream" && flags.Policy != "litellm" {
		return errors.New("--policy must be browser, upstream, or litellm")
	}
	return nil
}

type routeComparable struct {
	Rule     string
	Priority int
	Auth     nsl.AuthPolicy
}

func policyRoutes(name string, node nsl.Node, domain, policy string) []nsl.RouteInput {
	hostname := fmt.Sprintf("%s--%s.%s", strings.ToLower(name), node.Slug, domain)
	switch policy {
	case "browser":
		return []nsl.RouteInput{{Rule: fmt.Sprintf("Host(`%s`)", hostname), Priority: 100, Auth: nsl.AuthBrowser}}
	case "upstream":
		return []nsl.RouteInput{{Rule: fmt.Sprintf("Host(`%s`)", hostname), Priority: 100, Auth: nsl.AuthUpstream}}
	default:
		return []nsl.RouteInput{
			{Rule: fmt.Sprintf("Host(`%s`) && !(Path(`/v1`) || PathPrefix(`/v1/`))", hostname), Priority: 100, Auth: nsl.AuthBrowser},
			{Rule: fmt.Sprintf("Host(`%s`) && (Path(`/v1`) || PathPrefix(`/v1/`))", hostname), Priority: 110, Auth: nsl.AuthUpstream},
		}
	}
}
