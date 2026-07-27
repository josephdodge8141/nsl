package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	sub := os.Args[1]

	switch sub {
	case "list":
		fs := flag.NewFlagSet("list", flag.ExitOnError)
		apiURL := fs.String("api-url", defaultAPIURL(), "Registry API URL")
		fs.Parse(os.Args[2:])
		listCmd(*apiURL)

	case "add":
		fs := flag.NewFlagSet("add", flag.ExitOnError)
		apiURL := fs.String("api-url", defaultAPIURL(), "Registry API URL")
		var name, appType, targetURL, docsURL, connStr, desc string
		var noAuth, disabled bool
		fs.StringVar(&name, "name", "", "App name")
		fs.StringVar(&name, "n", "", "App name (shorthand)")
		fs.StringVar(&appType, "type", "", "App type (fe, be, db)")
		fs.StringVar(&appType, "t", "", "App type (shorthand)")
		fs.StringVar(&targetURL, "target-url", "", "Target URL")
		fs.StringVar(&targetURL, "u", "", "Target URL (shorthand)")
		fs.StringVar(&docsURL, "docs-url", "", "Docs URL")
		fs.StringVar(&connStr, "connection-string", "", "Postgres connection string")
		fs.StringVar(&desc, "description", "", "Description")
		fs.StringVar(&desc, "d", "", "Description (shorthand)")
		fs.BoolVar(&noAuth, "no-auth", false, "Disable auth")
		fs.BoolVar(&disabled, "disabled", false, "Create disabled")
		fs.Parse(os.Args[2:])

		addCmd(*apiURL, addFlags{
			Name:             name,
			AppType:          appType,
			TargetURL:        targetURL,
			DocsURL:          docsURL,
			ConnectionString: connStr,
			Description:      desc,
			NoAuth:           noAuth,
			Disabled:         disabled,
		})

	case "remove":
		fs := flag.NewFlagSet("remove", flag.ExitOnError)
		apiURL := fs.String("api-url", defaultAPIURL(), "Registry API URL")
		fs.Parse(os.Args[2:])
		removeCmd(*apiURL, fs.Args())

	case "help":
		usage()

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", sub)
		usage()
		os.Exit(1)
	}
}

func defaultAPIURL() string {
	if v := os.Getenv("NSL_API_URL"); v != "" {
		return v
	}
	return "http://localhost:7272"
}

func usage() {
	fmt.Fprint(os.Stderr, `Usage: nsl <command> [options]

Commands:
  list              List all apps
  add               Add a new app
  remove [id/name]  Remove an app

Options:
  --api-url         Registry API URL (default: http://localhost:7272, env: NSL_API_URL)

Run 'nsl help' for more information.
`)
}
