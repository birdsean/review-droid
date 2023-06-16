package entities

import (
	"fmt"
	"log"
	"os"

	"github.com/birdsean/review-droid/clients"
	"github.com/birdsean/review-droid/transformer"
	"github.com/google/go-github/v53/github"
)

var DEBUG = os.Getenv("DEBUG") == "true"

func Review(pr *github.PullRequest, client clients.GithubRepoClient) []*Comment {
	fmt.Printf("PR #%d: %s\n", pr.GetNumber(), pr.GetTitle())

	diff, err := client.GetPrDiff(pr)
	if err != nil {
		log.Fatalf("Failed to get raw diff: %v", err)
	}

	diffTransformer := transformer.DiffTransformer{}
	diffTransformer.Transform(diff)

	fileSegments := diffTransformer.GetFileSegments()
	allComments := []*Comment{}

	fmt.Printf("Getting comments for %d segments\n", len(fileSegments))
	for filename, segments := range fileSegments {
		for _, segment := range segments {
			comment := generateComments(segment)
			if comment == nil {
				continue
			}
			zippedComments, err := ZipComment(segment, *comment, filename, DEBUG)
			if err != nil {
				log.Fatalf("Failed to zip comment: %v", err)
			}
			allComments = append(allComments, zippedComments...)
			if DEBUG && len(allComments) >= 1 {
				break
			}
		}
		if DEBUG && len(allComments) >= 1 {
			break
		}
	}

	return allComments
}

func generateComments(segment string) *string {
	openAiClient := clients.NewOpenAiClient()
	completion, err := openAiClient.GetCompletion(segment, DEBUG)
	if err != nil {
		fmt.Printf("Failed to get completion: %v\n", err)
	}

	if DEBUG && completion != nil {
		fmt.Println("********************")
		fmt.Println(*completion)
		fmt.Println("********************")
	}
	return completion
}
