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
	openai_api "github.com/sashabaranov/go-openai"
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
		evaluateComments(pr, client)
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
			comment := generateComments(segment)
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

func generateComments(segment string) *string {
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

func evaluateComments(pr *github.PullRequest, client github_client.GithubRepoClient) error {
	// List review comments on a pull request
	comments, err := client.GetPrComments(pr)
	if err != nil {
		log.Fatalf("Failed to get comments: %v", err)
	}

	for _, commentDetails := range comments {
		comment := commentDetails.GetBody()
		diffHunk := commentDetails.GetDiffHunk()

		openaiClient := openai.OpenAiClient{}
		openaiClient.Init()
		conversation := []openai_api.ChatCompletionMessage{
			{
				Role:    openai_api.ChatMessageRoleSystem,
				Content: `You are an expert code review quality evaluator.`,
			},
			{
				Role:    openai_api.ChatMessageRoleUser,
				Content: fmt.Sprintf("Please summarize what this code is doing:\n%s", diffHunk),
			},
		}
		completion, err := openaiClient.RequestCompletion(conversation)
		if err != nil {
			return err
		}
		conversation = append(conversation, openai_api.ChatCompletionMessage{
			Role:    openai_api.ChatMessageRoleAssistant,
			Content: completion,
		}, openai_api.ChatCompletionMessage{
			Role:    openai_api.ChatMessageRoleUser,
			Content: fmt.Sprintf("Please rate the quality of this code review comment on a scale of 1-5:\n%s\nOnly respond with the number.", comment),
		})
		final, err := openaiClient.RequestCompletion(conversation)

		fmt.Printf("Comment: '%s' got a score of %s\n", comment, final)
	}
	return nil
}
