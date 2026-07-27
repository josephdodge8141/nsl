package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/josephdodge8141/nsl"
)

func listCmd(apiURL string) {
	client := nsl.NewClient(apiURL)
	apps, err := client.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tType\tRoute/Host\tEnabled")
	for _, a := range apps {
		id := a.ID
		if len(id) > 8 {
			id = id[:8]
		}
		route := a.RouteRule
		if route == "" {
			route = a.ContainerName
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\n", id, a.Name, a.AppType, route, a.Enabled)
	}
	w.Flush()
}
