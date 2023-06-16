package main

import (
	"log"
	"os"

	"github.com/birdsean/review-droid/clients"
	"github.com/birdsean/review-droid/entities"
)

var DEBUG = os.Getenv("DEBUG") == "true"

func main() {
	githubClient := clients.NewGithubRepoClient()

	prs, err := entities.GetAllPullRequests(githubClient)
	if err != nil {
		log.Fatalf("Failed to get pull requests: %v", err)
	}

	// Iterate over each pull request
	for _, pr := range prs {
		comments := pr.Review()
		err := pr.Comment(comments)
		if err != nil {
			log.Fatalf("Failed to comment: %v", err)
		}

		evalErr := pr.Evaluate()
		if evalErr != nil {
			log.Printf("Failed to evaluate comments: %v\n\n", err)
		}
	}
}
