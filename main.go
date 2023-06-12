package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	// Set up authentication using a personal access token
	token := os.Getenv("REVIEW_DROID_TOKEN")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	// Create a new GitHub client using the authentication client
	client := github.NewClient(tc)

	// Specify the repository details
	owner := os.Getenv("GITHUB_OWNER")
	repo := os.Getenv("GITHUB_REPO")

	// List pull requests for the specified repository
	prs, _, err := client.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
		State: "closed",
	})
	if err != nil {
		log.Fatalf("Failed to list pull requests: %v", err)
	}

	// Iterate over each pull request
	for _, pr := range prs {
		fmt.Printf("PR #%d: %s\n", pr.GetNumber(), pr.GetTitle())

		fmt.Printf("Getting diff for PR #%d\n", pr.GetNumber())
		diff, _, err := client.PullRequests.GetRaw(ctx, owner, repo, pr.GetNumber(), github.RawOptions{Type: github.Diff})
		if err != nil {
			log.Fatalf("Failed to get raw diff: %v", err)
		}
		fmt.Printf("Diff: %s\n", diff)
	}
}
