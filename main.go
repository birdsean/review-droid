package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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
				// TODO post comment to file.
			}
		}
		err := EvaluateReviewQuality(pr, client)
		if err != nil {
			log.Printf("Failed to evaluate comments: %v\n\n", err)
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

func EvaluateReviewQuality(pr *github.PullRequest, client github_client.GithubRepoClient) error {
	// List review comments on a pull request
	comments, err := client.GetPrComments(pr)
	if err != nil {
		log.Fatalf("Failed to get comments: %v", err)
	}

	// filter out comments that have been resolved
	filteredComments := []*github.PullRequestComment{}
	for _, comment := range comments {
		if comment.GetPosition() != 0 && comment.GetInReplyTo() == 0 {
			filteredComments = append(filteredComments, comment)
		}
	}

	fmt.Printf("Evaluating %d comments. %d were filtered.\n", len(filteredComments), len(comments)-len(filteredComments))
	return nil
	for _, commentDetails := range filteredComments {

		comment := commentDetails.GetBody()
		diffHunk := commentDetails.GetDiffHunk()

		openaiClient := openai.OpenAiClient{}
		openaiClient.Init()
		conversation := []openai_api.ChatCompletionMessage{
			{
				Role:    openai_api.ChatMessageRoleSystem,
				Content: `You are an expert code review quality evaluator. You rate code review comments on how practical, actionable, and pleasing to a programmer they are. ALL COMMENTS ABOUT IMPORTS GET A 1`,
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
		if err != nil {
			return err
		}

		// convert "final" first char to int
		firstChar := final[0:1]
		if strings.Contains("12345", firstChar) {
			score, err := strconv.Atoi(firstChar)
			if err != nil {
				return err
			}
			if score < 3 {
				fmt.Printf("Comment: '%s' got a score of %s. Deleting.\n", comment, final)
				err := client.DeleteComment(commentDetails)
				if err != nil {
					return err
				}
				continue
			} else if len(final) > 1 {
				// reply to the comment
				fmt.Printf("Comment: '%s' got a score of %s. Keeping.\n", comment, final)
				err := client.ReplyToComment(pr, commentDetails, fmt.Sprintf("Thank you for your feedback! The AI gives your comment a score of %s", final))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
