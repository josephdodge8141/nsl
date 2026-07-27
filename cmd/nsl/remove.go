package main

import (
	"fmt"
	"os"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/josephdodge8141/nsl"
)

func removeCmd(apiURL string, args []string) {
	client := nsl.NewClient(apiURL)

	if len(args) > 0 {
		idOrName := args[0]
		apps, err := client.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		var target *nsl.App
		for i, a := range apps {
			if a.ID == idOrName || strings.EqualFold(a.ID, idOrName) {
				target = &apps[i]
				break
			}
		}
		if target == nil {
			for i, a := range apps {
				if strings.HasPrefix(strings.ToLower(a.Name), strings.ToLower(idOrName)) {
					target = &apps[i]
					break
				}
			}
		}
		if target == nil {
			fmt.Fprintf(os.Stderr, "error: no app matching %q\n", idOrName)
			os.Exit(1)
		}

		if err := client.Delete(target.ID); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Removed %s\n", target.Name)
		return
	}

	apps, err := client.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if len(apps) == 0 {
		fmt.Println("No apps to remove.")
		return
	}

	opts := make([]string, len(apps))
	for i, a := range apps {
		label := a.Name
		if a.RouteRule != "" {
			label += " (" + a.RouteRule + ")"
		}
		opts[i] = label
	}

	var selected string
	if err := survey.AskOne(&survey.Select{Message: "Select app to remove:", Options: opts}, &selected); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var target *nsl.App
	for i, a := range apps {
		label := a.Name
		if a.RouteRule != "" {
			label += " (" + a.RouteRule + ")"
		}
		if label == selected {
			target = &apps[i]
			break
		}
	}

	if target == nil {
		fmt.Fprintf(os.Stderr, "error: unexpected selection\n")
		os.Exit(1)
	}

	confirm := false
	if err := survey.AskOne(&survey.Confirm{Message: fmt.Sprintf("Remove %s?", target.Name), Default: false}, &confirm); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if !confirm {
		fmt.Println("Cancelled.")
		return
	}

	if err := client.Delete(target.ID); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Removed %s\n", target.Name)
}
