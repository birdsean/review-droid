package main

import (
	"fmt"
	"log"

	"github.com/birdsean/review-droid/github"
)

func main() {
	client := github.GithubRepoClient{}
	client.Init()

	prs, err := client.GetPrs()
	if err != nil {
		log.Fatalf("Failed to get pull requests: %v", err)
	}

	// Iterate over each pull request
	for _, pr := range prs {
		fmt.Printf("PR #%d: %s\n", pr.GetNumber(), pr.GetTitle())

		fmt.Printf("Getting diff for PR #%d\n", pr.GetNumber())
		diff, err := client.GetPrDiff(pr)
		if err != nil {
			log.Fatalf("Failed to get raw diff: %v", err)
		}
		fmt.Printf("Diff: %s\n", diff)
	}
}
