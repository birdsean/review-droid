package main

import (
	"fmt"
	"log"

	"github.com/birdsean/review-droid/comments"
	"github.com/birdsean/review-droid/github"
	"github.com/birdsean/review-droid/openai"
	"github.com/birdsean/review-droid/transformer"
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

		diffTransformer := transformer.DiffTransformer{}
		diffTransformer.Transform(diff)

		fileSegments := diffTransformer.GetFileSegments()
		allComments := []*comments.Comment{}
		fmt.Printf("Getting comments for %d segments\n", len(fileSegments))
		for filename, segments := range fileSegments {
			fmt.Print("///////////////////////////////////////////\n")
			fmt.Printf("Getting comments for file: %s\n", filename)
			fmt.Printf("Getting comments for segment: %s\n", segments)
			// comment := retrieveComments(segment)
			// zippedComments, err := comments.ZipComment(segment, comment)
			// if err != nil {
			// 	log.Fatalf("Failed to zip comment: %v", err)
			// }
			// allComments = append(allComments, zippedComments...)
		}

		fmt.Printf("allComments: %v\n", allComments)
	}
}

func retrieveComments(segment string) string {
	openAiClient := openai.OpenAiClient{}
	openAiClient.Init()
	completion, err := openAiClient.GetCompletion(segment)
	if err != nil {
		fmt.Printf("Failed to get completion: %v\n", err)
	}

	fmt.Println("********************")
	fmt.Println(*completion)
	fmt.Println("********************")
	return *completion
}
