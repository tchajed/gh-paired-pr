package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v81/github"
	"github.com/tchajed/gh-paired-pr/check_pr"
)

var verbose bool
var showCommit bool

func verbosePrintf(format string, args ...any) {
	if verbose {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

func normalizeRepoToSlug(url string) string {
	url = strings.TrimPrefix(url, "https://github.com/")
	url = strings.TrimPrefix(url, "github.com/")
	return url
}

func main() {
	var (
		baseRepo      string
		prNum         int
		dependentRepo string
	)

	flag.StringVar(&baseRepo, "base", "", "base repo to check a PR from")
	flag.StringVar(&dependentRepo, "dependency", "", "repo to look for a dependent PR")
	flag.IntVar(&prNum, "pr", 0, "PR to check in base repo")
	flag.BoolVar(&verbose, "verbose", false, "print extra info to stderr")
	flag.BoolVar(&showCommit, "commit", false, "include commit hash in output (URL<tab>hash)")
	flag.Parse()

	baseRepo = normalizeRepoToSlug(baseRepo)
	dependentRepo = normalizeRepoToSlug(dependentRepo)

	token := os.Getenv("GITHUB_TOKEN")

	if prNum == 0 {
		fmt.Fprintf(os.Stderr, "No PR specified\n")
		os.Exit(1)
	}
	if baseRepo == "" {
		fmt.Fprintf(os.Stderr, "No base repo specified\n")
		os.Exit(1)
	}
	if dependentRepo == "" {
		fmt.Fprintf(os.Stderr, "No dependent repo specified\n")
		os.Exit(1)
	}

	client := github.NewClient(nil)
	if token != "" {
		client = client.WithAuthToken(token)
	}
	info, err := check_pr.CheckPrDependency(client, baseRepo, prNum, dependentRepo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking PR dependency: %v\n", err)
		os.Exit(1)
	}

	if !info.HasDependency {
		verbosePrintf("no dependent PR from %s found\n", dependentRepo)
		return
	}

	verbosePrintf("INFO depends on: %s#%d\n", info.DependentSlug, info.DependentNum)
	verbosePrintf("INFO status: %s\n", info.DependentPr.GetState())
	verbosePrintf("INFO source: %s at %s (%s)\n", info.SourceSlug, info.SourceRef, info.SourceSHA)
	if !(info.DependentPr.GetState() == "open") {
		return
	}
	if showCommit {
		fmt.Printf("%s#%s\n", info.SourceUrl(), info.SourceSHA)
	} else {
		fmt.Printf("%s\n", info.SourceUrl())
	}
}
