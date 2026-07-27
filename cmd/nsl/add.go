package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/josephdodge8141/nsl"
)

type addFlags struct {
	Name             string
	AppType          string
	TargetURL        string
	DocsURL          string
	ConnectionString string
	Description      string
	NoAuth           bool
	Disabled         bool
}

var hostnameRE = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?$`)

func addCmd(apiURL string, f addFlags) {
	client := nsl.NewClient(apiURL)
	interactive := f.Name == "" && f.AppType == ""

	if interactive {
		interactiveAdd(client, f)
		return
	}

	if err := validateAddFlags(f); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if f.Name == "" {
		f.Name = promptName(client)
	}
	if f.AppType == "" {
		f.AppType = promptType()
	}
	if f.Description == "" {
		var desc string
		if err := survey.AskOne(&survey.Input{Message: "Description (optional):"}, &desc); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		f.Description = desc
	}

	conditionalFields(client, &f, apiURL)

	showSummary(f)
	confirm := false
	if err := survey.AskOne(&survey.Confirm{Message: "Create this app?", Default: true}, &confirm); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if !confirm {
		fmt.Println("Cancelled.")
		return
	}

	app := nsl.App{
		Name:             f.Name,
		Description:      f.Description,
		AppType:          f.AppType,
		TargetURL:        f.TargetURL,
		DocsURL:          f.DocsURL,
		ConnectionString: f.ConnectionString,
		NoAuth:           f.NoAuth,
		Enabled:          !f.Disabled,
	}

	created, err := client.Create(app)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created %s -> %s\n", created.Name, created.RouteRule)
}

func validateAddFlags(f addFlags) error {
	if f.Name != "" && !hostnameRE.MatchString(f.Name) {
		return errors.New("name must be a valid hostname: alphanumeric and hyphens only")
	}
	if f.AppType != "" && f.AppType != "fe" && f.AppType != "be" && f.AppType != "db" {
		return errors.New("type must be one of: fe, be, db")
	}
	if f.TargetURL != "" {
		if _, err := url.ParseRequestURI(f.TargetURL); err != nil {
			return fmt.Errorf("invalid target URL: %w", err)
		}
	}
	if f.DocsURL != "" {
		if _, err := url.ParseRequestURI(f.DocsURL); err != nil {
			return fmt.Errorf("invalid docs URL: %w", err)
		}
	}
	if f.ConnectionString != "" && !strings.HasPrefix(f.ConnectionString, "postgres://") {
		return errors.New("connection string must start with postgres://")
	}
	if f.AppType == "fe" && f.TargetURL == "" {
		return errors.New("target-url is required for type fe")
	}
	if f.AppType == "be" && f.DocsURL == "" {
		return errors.New("docs-url is required for type be")
	}
	if f.AppType == "db" && f.ConnectionString == "" {
		return errors.New("connection-string is required for type db")
	}
	return nil
}

func promptName(client *nsl.Client) string {
	var name string
	existing := fetchNames(client)
	validate := func(val interface{}) error {
		s, ok := val.(string)
		if !ok {
			return errors.New("invalid input")
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return errors.New("name is required")
		}
		if !hostnameRE.MatchString(s) {
			return errors.New("name must be a valid hostname: alphanumeric and hyphens only")
		}
		for _, n := range existing {
			if strings.EqualFold(s, n) {
				return fmt.Errorf("name %q already exists", s)
			}
		}
		return nil
	}
	if err := survey.AskOne(&survey.Input{Message: "Name:"}, &name, survey.WithValidator(validate)); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	return name
}

func promptType() string {
	var t string
	if err := survey.AskOne(&survey.Select{Message: "Type:", Options: []string{"fe", "be", "db"}}, &t); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	return t
}

func promptURL(msg string, required bool) string {
	var u string
	validate := func(val interface{}) error {
		s, ok := val.(string)
		if !ok {
			return errors.New("invalid input")
		}
		s = strings.TrimSpace(s)
		if !required && s == "" {
			return nil
		}
		if s == "" {
			return fmt.Errorf("%s is required", msg)
		}
		if _, err := url.ParseRequestURI(s); err != nil {
			return fmt.Errorf("valid URL required: %w", err)
		}
		return nil
	}
	if err := survey.AskOne(&survey.Input{Message: msg + ":"}, &u, survey.WithValidator(validate)); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	return u
}

func promptConnString() string {
	var s string
	validate := func(val interface{}) error {
		v, ok := val.(string)
		if !ok {
			return errors.New("invalid input")
		}
		v = strings.TrimSpace(v)
		if v == "" {
			return errors.New("connection string is required")
		}
		if !strings.HasPrefix(v, "postgres://") {
			return errors.New("connection string must start with postgres://")
		}
		return nil
	}
	if err := survey.AskOne(&survey.Input{Message: "Connection string:"}, &s, survey.WithValidator(validate)); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	return s
}

func conditionalFields(client *nsl.Client, f *addFlags, apiURL string) {
	switch f.AppType {
	case "fe":
		if f.TargetURL == "" {
			f.TargetURL = promptURL("Target URL", true)
		}
	case "be":
		if f.DocsURL == "" {
			f.DocsURL = promptURL("Docs URL", true)
		}
		if f.TargetURL == "" {
			f.TargetURL = promptURL("Target URL (optional)", false)
		}
	case "db":
		if f.ConnectionString == "" {
			f.ConnectionString = promptConnString()
		}
	}

	if !f.NoAuth {
		noAuth := false
		survey.AskOne(&survey.Confirm{Message: "No auth?", Default: false}, &noAuth)
		f.NoAuth = noAuth
	}

	if !f.Disabled {
		disabled := false
		survey.AskOne(&survey.Confirm{Message: "Disabled?", Default: false}, &disabled)
		f.Disabled = disabled
	}
}

func interactiveAdd(client *nsl.Client, f addFlags) {
	f.Name = promptName(client)
	f.AppType = promptType()

	var desc string
	survey.AskOne(&survey.Input{Message: "Description (optional):"}, &desc)
	f.Description = desc

	conditionalFields(client, &f, client.BaseURL)

	showSummary(f)
	confirm := false
	survey.AskOne(&survey.Confirm{Message: "Create this app?", Default: true}, &confirm)
	if !confirm {
		fmt.Println("Cancelled.")
		return
	}

	app := nsl.App{
		Name:             f.Name,
		Description:      f.Description,
		AppType:          f.AppType,
		TargetURL:        f.TargetURL,
		DocsURL:          f.DocsURL,
		ConnectionString: f.ConnectionString,
		NoAuth:           f.NoAuth,
		Enabled:          !f.Disabled,
	}

	created, err := client.Create(app)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created %s -> %s\n", created.Name, created.RouteRule)
}

func showSummary(f addFlags) {
	fmt.Println("\n--- Summary ---")
	fmt.Printf("  Name:        %s\n", f.Name)
	fmt.Printf("  Type:        %s\n", f.AppType)
	if f.TargetURL != "" {
		fmt.Printf("  Target URL:  %s\n", f.TargetURL)
	}
	if f.DocsURL != "" {
		fmt.Printf("  Docs URL:    %s\n", f.DocsURL)
	}
	if f.ConnectionString != "" {
		fmt.Printf("  Connection:  %s\n", f.ConnectionString)
	}
	if f.Description != "" {
		fmt.Printf("  Description: %s\n", f.Description)
	}
	fmt.Printf("  No auth:     %v\n", f.NoAuth)
	fmt.Printf("  Disabled:    %v\n", f.Disabled)
	fmt.Println()
}

func fetchNames(client *nsl.Client) []string {
	apps, err := client.List()
	if err != nil {
		return nil
	}
	names := make([]string, len(apps))
	for i, a := range apps {
		names[i] = a.Name
	}
	return names
}
