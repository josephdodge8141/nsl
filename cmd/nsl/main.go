// Package main implements the nsl command-line client.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/josephdodge8141/nsl"
)

// Version is embedded at build time for client/server compatibility reporting.
var Version = "dev"

type cliError struct {
	message string
	help    string
	usage   bool
}

func (e *cliError) Error() string { return e.message }

func main() {
	if err := run(os.Args[1:]); err != nil {
		var structured *cliError
		if errors.As(err, &structured) {
			fmt.Printf("error: %s\n", quote(structured.message))
			if structured.help != "" {
				fmt.Printf("help: %s\n", quote(structured.help))
			}
			if structured.usage {
				os.Exit(2)
			}
		} else {
			fmt.Printf("error: %s\n", quote(err.Error()))
		}
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printIdentity()
		return listCmd(defaultAPIURL())
	}
	switch args[0] {
	case "list":
		fs := commandFlags("list")
		apiURL := fs.String("api-url", defaultAPIURL(), "Registry API URL")
		if err := fs.Parse(args[1:]); errors.Is(err, flag.ErrHelp) {
			commandHelp("list")
			return nil
		} else if err != nil {
			return usageError(err.Error(), "nsl list [--api-url URL]")
		}
		return listCmd(*apiURL)
	case "nodes":
		fs := commandFlags("nodes")
		apiURL := fs.String("api-url", defaultAPIURL(), "Registry API URL")
		if err := fs.Parse(args[1:]); errors.Is(err, flag.ErrHelp) {
			commandHelp("nodes")
			return nil
		} else if err != nil {
			return usageError(err.Error(), "nsl nodes [--api-url URL]")
		}
		return nodesCmd(*apiURL)
	case "add":
		fs := commandFlags("add")
		apiURL := fs.String("api-url", defaultAPIURL(), "Registry API URL")
		flags := addFlags{Policy: "browser"}
		fs.StringVar(&flags.Name, "name", "", "App name (required)")
		fs.StringVar(&flags.TargetURL, "target-url", "", "HTTP target URL (required)")
		fs.StringVar(&flags.Description, "description", "", "Description")
		fs.StringVar(&flags.Policy, "policy", "browser", "browser, upstream, or litellm")
		fs.BoolVar(&flags.Disabled, "disabled", false, "Register disabled")
		if err := fs.Parse(args[1:]); errors.Is(err, flag.ErrHelp) {
			commandHelp("add")
			return nil
		} else if err != nil {
			return usageError(err.Error(), "nsl add --name <name> --target-url <url>")
		}
		return addCmd(*apiURL, flags)
	case "remove":
		fs := commandFlags("remove")
		apiURL := fs.String("api-url", defaultAPIURL(), "Registry API URL")
		if err := fs.Parse(args[1:]); errors.Is(err, flag.ErrHelp) {
			commandHelp("remove")
			return nil
		} else if err != nil {
			return usageError(err.Error(), "nsl remove <id-or-name>")
		}
		if fs.NArg() != 1 {
			return usageError("remove requires one app id or exact name", "nsl remove <id-or-name>")
		}
		return removeCmd(*apiURL, fs.Arg(0))
	case "version":
		if len(args) > 1 && (args[1] == "--help" || args[1] == "-h") {
			fmt.Println(`description: Show CLI and registry versions
usage: "nsl version"`)
			return nil
		}
		fmt.Printf("nsl: %s\n", quote(Version))
		serverVersion, err := nsl.NewClient(defaultAPIURL()).FetchVersion()
		if err != nil {
			return err
		}
		fmt.Printf("registry: %s\n", quote(serverVersion))
		return nil
	case "enrollment-token":
		fs := commandFlags("enrollment-token")
		nodeName := fs.String("node-name", "", "New node's canonical name (required)")
		brokerURL := fs.String("broker-url", defaultBrokerURL(), "Enrollment broker URL")
		ttl := fs.Duration("ttl", 10*time.Minute, "Token lifetime (maximum 15m)")
		if err := fs.Parse(args[1:]); errors.Is(err, flag.ErrHelp) {
			commandHelp("enrollment-token")
			return nil
		} else if err != nil {
			return usageError(err.Error(), "nsl enrollment-token --node-name <name>")
		}
		if *nodeName == "" {
			return usageError("--node-name is required", "nsl enrollment-token --node-name <name>")
		}
		if !regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,28}[a-z0-9])?$`).MatchString(*nodeName) {
			return usageError("--node-name must be a lowercase DNS label of at most 30 characters", "nsl enrollment-token --node-name laptop3")
		}
		if *ttl < time.Minute || *ttl > 15*time.Minute {
			return usageError("--ttl must be between 1m and 15m", "nsl enrollment-token --node-name <name> --ttl 10m")
		}
		adminToken := os.Getenv("NSL_BROKER_ADMIN_TOKEN")
		if adminToken == "" {
			return usageError("NSL_BROKER_ADMIN_TOKEN is required", "Set the broker admin token in the environment")
		}
		token, err := nsl.IssueEnrollmentToken(context.Background(), *brokerURL, adminToken, *nodeName, *ttl)
		if err != nil {
			return err
		}
		fmt.Println("enrollment:")
		fmt.Printf("  node_name: %s\n", quote(token.NodeName))
		fmt.Printf("  token: %s\n", quote(token.EnrollmentToken))
		fmt.Printf("  expires_at: %s\n", quote(token.ExpiresAt.Format(time.RFC3339)))
		fmt.Println(`help: "Set this once as NSL_ENROLLMENT_TOKEN on the new node"`)
		return nil
	case "help", "--help", "-h":
		usage()
		return nil
	default:
		return usageError("unknown command "+args[0], "nsl help")
	}
}

func commandFlags(name string) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	return flags
}

func defaultAPIURL() string {
	if value := os.Getenv("NSL_API_URL"); value != "" {
		return value
	}
	return "http://localhost:7272"
}

func defaultBrokerURL() string {
	if value := os.Getenv("NSL_BROKER_URL"); value != "" {
		return value
	}
	return "https://nsl-enrollment-broker.example.workers.dev"
}

func usageError(message, help string) error {
	return &cliError{message: message, help: help, usage: true}
}

func quote(value string) string { return strconv.Quote(value) }

func printIdentity() {
	executable, err := os.Executable()
	if err != nil {
		executable = "nsl"
	}
	home, _ := os.UserHomeDir()
	if home != "" {
		if relative, err := filepath.Rel(home, executable); err == nil && relative != ".." {
			executable = "~/" + relative
		}
	}
	fmt.Printf("bin: %s\n", quote(executable))
	fmt.Println("description: Manage applications across Not-So-Localhost machines")
}

func usage() {
	fmt.Println(`description: Manage applications across Not-So-Localhost machines
commands[6]{name,usage}:
  "list","nsl list [--api-url URL]"
  "nodes","nsl nodes [--api-url URL]"
  "add","nsl add --name <name> --target-url <url> [--policy browser|upstream|litellm]"
  "remove","nsl remove <id-or-name>"
  "version","nsl version"
  "enrollment-token","nsl enrollment-token --node-name <name>"`)
}

func commandHelp(command string) {
	switch command {
	case "add":
		fmt.Println(`description: Register a direct HTTP service
usage: "nsl add --name <name> --target-url <url> [--policy browser|upstream|litellm]"`)
	case "remove":
		fmt.Println(`description: Remove an app registration
usage: "nsl remove <id-or-exact-name>"`)
	case "enrollment-token":
		fmt.Println(`description: Issue a one-time node enrollment token
usage: "nsl enrollment-token --node-name <name> [--ttl 10m]"`)
	default:
		fmt.Printf("description: %s registered NSL state\nusage: %s\n", command, quote("nsl "+command+" [--api-url URL]"))
	}
}
