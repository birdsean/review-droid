package main

import (
	"fmt"
	"log"

	"github.com/birdsean/review-droid/comments"
	"github.com/birdsean/review-droid/github_client"
	"github.com/birdsean/review-droid/openai"
	"github.com/birdsean/review-droid/transformer"
	"github.com/google/go-github/github"
)

func main() {
	client := github_client.GithubRepoClient{}
	client.Init()

	prs, err := client.GetPrs()
	if err != nil {
		log.Fatalf("Failed to get pull requests: %v", err)
	}

	// Iterate over each pull request
	for _, pr := range prs {
		reviewPR(pr, client)
	}
}

func reviewPR(pr *github.PullRequest, client github_client.GithubRepoClient) {
	fmt.Printf("PR #%d: %s\n", pr.GetNumber(), pr.GetTitle())

	fmt.Printf("Getting diff for PR #%d\n", pr.GetNumber())
	diff, err := client.GetPrDiff(pr)
	if err != nil {
		log.Fatalf("Failed to get raw diff: %v", err)
	}

	diffTransformer := transformer.DiffTransformer{}
	diffTransformer.Transform(diff)

	fileSegments := diffTransformer.GetFileSegments()
	allComments := []*comments.Comment{}
	failedComments := []string{}

	fmt.Printf("Getting comments for %d segments\n", len(fileSegments))
	for filename, segments := range fileSegments {
		for _, segment := range segments {
			comment := retrieveComments(segment)
			if comment == nil {
				failedComments = append(failedComments, segment)
				continue
			}
			zippedComments, err := comments.ZipComment(segment, *comment, filename)
			if err != nil {
				log.Fatalf("Failed to zip comment: %v", err)
			}
			allComments = append(allComments, zippedComments...)
		}
	}

	fmt.Println("Comments:")
	for _, comment := range allComments {
		fmt.Printf("%+v\n", comment)
	}

	if len(failedComments) == 0 {
		fmt.Println("Failed to get comments for the following segments:")
		for _, comment := range failedComments {
			fmt.Printf("Failed to get comment for segment: %s\n", comment)
		}
	}
}

func retrieveComments(segment string) *string {
	openAiClient := openai.OpenAiClient{}
	openAiClient.Init()
	completion, err := openAiClient.GetCompletion(segment)
	if err != nil {
		fmt.Printf("Failed to get completion: %v\n", err)
	}

	if completion != nil {
		fmt.Println("********************")
		fmt.Println(*completion)
		fmt.Println("********************")
	}
	return completion
}
