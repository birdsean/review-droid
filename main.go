package main

import (
	"fmt"
	"log"
	"os"

	"github.com/birdsean/review-droid/comments"
	"github.com/birdsean/review-droid/github_client"
	"github.com/birdsean/review-droid/openai"
	"github.com/birdsean/review-droid/transformer"
	"github.com/google/go-github/v53/github"
)

var DEBUG = os.Getenv("DEBUG") == "true"

func main() {
	client := github_client.GithubRepoClient{}
	client.Init()

	prs, err := client.GetPrs()
	if err != nil {
		log.Fatalf("Failed to get pull requests: %v", err)
	}

	// Iterate over each pull request
	for _, pr := range prs {
		comments := reviewPR(pr, client)
		commitId := pr.GetHead().GetSHA()
		fmt.Printf("Posting %d comments to PR #%d\n", len(comments), pr.GetNumber())
		for _, comment := range comments {
			ghComment := client.ParsedCommentToGithubComment(comment, commitId)
			err := client.PostComment(pr, ghComment)
			if err != nil {
				fmt.Printf("Failed to post comment: %v\n", err)
				fmt.Printf("Comment: %v\n", ghComment)
			}
		}
	}
}

func reviewPR(pr *github.PullRequest, client github_client.GithubRepoClient) []*comments.Comment {
	fmt.Printf("PR #%d: %s\n", pr.GetNumber(), pr.GetTitle())

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
		for _, segment := range segments {
			comment := retrieveComments(segment)
			if comment == nil {
				continue
			}
			zippedComments, err := comments.ZipComment(segment, *comment, filename, DEBUG)
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

func retrieveComments(segment string) *string {
	openAiClient := openai.OpenAiClient{}
	openAiClient.Init()
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
