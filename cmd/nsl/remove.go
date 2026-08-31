// Package main implements the nsl command-line client.
package main

import (
	"fmt"
	"strings"

	"github.com/josephdodge8141/nsl"
)

func removeCmd(apiURL, identifier string) error {
	if identifier == "" {
		return usageError("app id or exact name is required", "nsl remove <id-or-name>")
	}
	client := nsl.NewClient(apiURL)
	node, err := client.LocalNode()
	if err != nil {
		return err
	}
	apps, err := client.List()
	if err != nil {
		return err
	}
	matches := make([]nsl.App, 0, 1)
	for _, app := range apps {
		if app.NodeID == node.ID && (app.ID == identifier || strings.EqualFold(app.Name, identifier)) {
			matches = append(matches, app)
		}
	}
	if len(matches) == 0 {
		fmt.Println("app:")
		fmt.Printf("  name: %s\n", quote(identifier))
		fmt.Println("  status: already_absent")
		return nil
	}
	if len(matches) > 1 {
		return usageError("name matches apps on multiple nodes; use the app id", "nsl list")
	}
	if err := client.Delete(matches[0]); err != nil {
		return err
	}
	fmt.Println("app:")
	fmt.Printf("  id: %s\n", quote(matches[0].ID))
	fmt.Printf("  name: %s\n", quote(matches[0].Name))
	fmt.Println("  status: removed")
	return nil
}
