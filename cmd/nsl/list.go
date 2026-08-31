package main

import (
	"fmt"

	"github.com/josephdodge8141/nsl"
)

func listCmd(apiURL string) error {
	apps, err := nsl.NewClient(apiURL).List()
	if err != nil {
		return err
	}
	if len(apps) == 0 {
		fmt.Println("apps: 0 registered apps")
		fmt.Println(`help[1]: "Run nsl add --name <name> --target-url <url> to register one"`)
		return nil
	}
	fmt.Printf("count: %d total\n", len(apps))
	fmt.Printf("apps[%d]{id,name,node_id,url}:\n", len(apps))
	for _, app := range apps {
		fmt.Printf("  %s,%s,%s,%s\n", quote(app.ID), quote(app.Name), quote(app.NodeID), quote(app.PublicURL))
	}
	fmt.Println(`help[2]: "Run nsl nodes to list machine identities","Run nsl remove <id-or-name> to unregister an app"`)
	return nil
}

func nodesCmd(apiURL string) error {
	nodes, err := nsl.NewClient(apiURL).Nodes()
	if err != nil {
		return err
	}
	fmt.Printf("count: %d total\n", len(nodes))
	fmt.Printf("nodes[%d]{id,name,slug,roles}:\n", len(nodes))
	for _, node := range nodes {
		fmt.Printf("  %s,%s,%s,%s\n", quote(node.ID), quote(node.Name), quote(node.Slug), quote(joinRoles(node.Roles)))
	}
	return nil
}

func joinRoles(roles []string) string {
	if len(roles) == 0 {
		return "none"
	}
	result := roles[0]
	for _, role := range roles[1:] {
		result += "+" + role
	}
	return result
}
