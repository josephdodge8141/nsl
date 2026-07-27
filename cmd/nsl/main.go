package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/josephdodge8141/nsl"
)

var Version = "dev"

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
		noVersionCheck := fs.Bool("no-version-check", false, "Skip server version check")
		fs.Parse(os.Args[2:])
		maybeCheckVersion(*apiURL, *noVersionCheck)
		listCmd(*apiURL)

	case "add":
		fs := flag.NewFlagSet("add", flag.ExitOnError)
		apiURL := fs.String("api-url", defaultAPIURL(), "Registry API URL")
		noVersionCheck := fs.Bool("no-version-check", false, "Skip server version check")
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

		maybeCheckVersion(*apiURL, *noVersionCheck)
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
		noVersionCheck := fs.Bool("no-version-check", false, "Skip server version check")
		fs.Parse(os.Args[2:])
		maybeCheckVersion(*apiURL, *noVersionCheck)
		removeCmd(*apiURL, fs.Args())

	case "version":
		fmt.Printf("nsl %s\n", Version)
		client := nsl.NewClient(defaultAPIURL())
		if sv, err := client.FetchVersion(); err == nil {
			fmt.Printf("registry %s\n", sv)
		} else {
			fmt.Fprintf(os.Stderr, "registry: %v\n", err)
		}

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

func maybeCheckVersion(apiURL string, skip bool) {
	if skip {
		return
	}
	client := nsl.NewClient(apiURL)
	sv, err := client.FetchVersion()
	if err != nil {
		return
	}
	if sv != Version {
		fmt.Fprintf(os.Stderr, "warning: nsl %s, registry %s (use --no-version-check to suppress)\n", Version, sv)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `Usage: nsl <command> [options]

Commands:
  list              List all apps
  add               Add a new app
  remove [id/name]  Remove an app
  version           Show version info

Options:
  --api-url         Registry API URL (default: http://localhost:7272, env: NSL_API_URL)

Run 'nsl help' for more information.
` + "\n")
}
