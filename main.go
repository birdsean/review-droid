package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

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

		segments := diffTransformer.GetSegments()
		allComments := []string{}
		fmt.Printf("Getting commments for %d segments\n", len(segments))
		for _, segment := range segments {
			comment := retrieveComments(segment)
			allComments = append(allComments, comment)
		}

		writeResults(allComments)
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

func writeResults(comments []string) {
	fileContents, err := json.Marshal(comments)
	if err != nil {
		log.Fatalf("Failed to marshal comments: %v", err)
	}
	// save fileContents to results.json
	err = ioutil.WriteFile("results.json.test", fileContents, 0644)
	if err != nil {
		log.Fatalf("Failed to write comments to file: %v", err)
	}
}
