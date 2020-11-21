// CLI tool to merge PRs with specified labels that have passed checks.
//
// This tool is intended to be run at a regular interval (e.g. using GitHub
// workflows).
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v32/github"
)

var (
	tokenFlag = flag.String(
		"token",
		os.Getenv("GITHUB_TOKEN"),
		"GitHub token used for authentication. Uses GITHUB_TOKEN if not provided.",
	)
	repoFlag = flag.String(
		"repository",
		os.Getenv("GITHUB_REPOSITORY"),
		"GitHub repository to check issues on. Should be of the for <owner>/<repo>. Uses GITHUB_REPOSITORY if not provided.",
	)
	labelFlag = flag.String(
		"label",
		"",
		"Label to filter pull requests by. Only PRs with this label will be checked and merged.",
	)
)

func init() {
	flag.Parse()
}

func main() {
	token := *tokenFlag
	if strings.TrimSpace(token) == "" {
		log.Fatal("GitHub token not provided via CLI or environment variable.")
	}

	repo := *repoFlag
	if strings.TrimSpace(repo) == "" {
		log.Fatal("GitHub repository not provided via CLI or environment variable.")
	}

	label := *labelFlag
	if strings.TrimSpace(label) == "" {
		log.Fatal("Label filter not provided.")
	}

	repoParts := strings.Split(repo, "/")
	if len(repoParts) != 2 {
		log.Fatalf("Expected GitHub repository name to be of the form <owner>/<repo>. '%s' is not.", repo)
	}
	owner := repoParts[0]
	repoName := repoParts[1]

	client := github.NewClient(nil)
	pullRequests, _, err := client.PullRequests.List(context.TODO(), owner, repoName, &github.PullRequestListOptions{})
	if err != nil {
		log.Fatalf("Failed to retrieve pull requests from %s: %v", repo, err)
	}
	log.Printf("Retrieved a total of %d pull requests from %s", len(pullRequests), repo)

	labeledPullRequests := filterPullRequestsByLabel(pullRequests, label)
	log.Printf("Found %d pull requests in %s with the label %s", len(labeledPullRequests), repo, label)

	failureCount := 0
	for _, pullRequest := range labeledPullRequests {
		if err := checkAndMerge(context.TODO(), client, owner, repoName, pullRequest); err != nil {
			log.Print(err)
			failureCount++
		}
	}

	if failureCount > 0 {
		log.Fatalf(
			"Failed to check and merge %d/%d pull requests. See the above logs for details.",
			failureCount,
			len(labeledPullRequests),
		)
	}
}

func filterPullRequestsByLabel(pullRequests []*github.PullRequest, expectedLabel string) []*github.PullRequest {
	filteredPullRequests := []*github.PullRequest{}
	for _, pullRequest := range pullRequests {
		contains := false
		for _, label := range pullRequest.Labels {
			if label.GetName() == expectedLabel {
				contains = true
			}
		}
		if contains {
			filteredPullRequests = append(filteredPullRequests, pullRequest)
		}
	}
	return filteredPullRequests
}

func checkAndMerge(ctx context.Context, client *github.Client, owner, repoName string, pullRequest *github.PullRequest) error {
	if pullRequest.GetMergeable() {
		result, _, err := client.PullRequests.Merge(
			ctx,
			owner,
			repoName,
			pullRequest.GetNumber(),
			"Merged by merged",
			&github.PullRequestOptions{},
		)
		if err != nil {
			return fmt.Errorf("failed to merge mergeable pull request %d: %w", pullRequest.GetNumber(), err)
		}
		log.Printf(
			"Merged pull request %d into %s as commit %s",
			pullRequest.GetNumber(),
			pullRequest.GetBase().GetLabel(),
			result.GetSHA(),
		)
	} else {
		log.Printf("Pull request %d is not mergeable. Skipping.", pullRequest.GetNumber())
	}
	return nil
}
