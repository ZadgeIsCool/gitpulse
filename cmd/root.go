package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ZadgeIsCool/gitpulse/analysis"
	"github.com/ZadgeIsCool/gitpulse/dashboard"
	"github.com/ZadgeIsCool/gitpulse/github"
)

const version = "0.1.0"

func Execute() error {
	var (
		token      string
		days       int
		showVer    bool
		jsonOutput bool
	)

	flag.StringVar(&token, "token", "", "GitHub personal access token (or set GITHUB_TOKEN env var)")
	flag.IntVar(&days, "days", 90, "Number of days to analyze (default: 90)")
	flag.BoolVar(&showVer, "version", false, "Print version and exit")
	flag.BoolVar(&jsonOutput, "json", false, "Output as JSON instead of dashboard")
	flag.Usage = usage
	flag.Parse()

	if showVer {
		fmt.Printf("gitpulse v%s\n", version)
		return nil
	}

	if flag.NArg() < 1 {
		usage()
		return errors.New("repository argument required (owner/repo)")
	}

	repo := flag.Arg(0)
	parts := strings.SplitN(repo, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("invalid repository format %q, expected owner/repo", repo)
	}
	owner, name := parts[0], parts[1]

	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		return errors.New("GitHub token required: use --token flag or set GITHUB_TOKEN env var")
	}

	client := github.NewClient(token)

	fmt.Printf("Analyzing %s/%s (last %d days)...\n\n", owner, name, days)

	report, err := analysis.GenerateReport(client, owner, name, days)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	if jsonOutput {
		return dashboard.PrintJSON(report)
	}
	dashboard.Print(report)
	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, `gitpulse v%s — health check dashboard for GitHub repos

Usage:
  gitpulse [flags] owner/repo

Examples:
  gitpulse golang/go
  gitpulse --days 30 rust-lang/rust
  gitpulse --json kubernetes/kubernetes

Flags:
`, version)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, `
Environment:
  GITHUB_TOKEN   GitHub personal access token (alternative to --token)`)
}
