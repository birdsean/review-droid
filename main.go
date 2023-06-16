package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/birdsean/review-droid/clients"
	"github.com/birdsean/review-droid/entities"
	"github.com/google/go-github/v53/github"
	openai_api "github.com/sashabaranov/go-openai"
)

var DEBUG = os.Getenv("DEBUG") == "true"

func main() {
	githubClient := clients.NewGithubRepoClient()

	prs, err := githubClient.GetPrs()
	if err != nil {
		log.Fatalf("Failed to get pull requests: %v", err)
	}

	// Iterate over each pull request
	for _, pr := range prs {
		comments := entities.Review(pr, githubClient) // TODO: call class instead of function
		commitId := pr.GetHead().GetSHA()
		fmt.Printf("Posting %d comments to PR #%d\n", len(comments), pr.GetNumber())
		for _, comment := range comments {
			ghComment := entities.ParsedCommentToGithubComment(comment, commitId)
			err := githubClient.PostComment(pr, ghComment)
			if err != nil {
				fmt.Printf("Failed to post comment: %v\n", err)
				// TODO post comment to entire file if failed on a line.
			}
		}
		err := EvaluateReviewQuality(pr, githubClient)
		if err != nil {
			log.Printf("Failed to evaluate comments: %v\n\n", err)
		}
	}
}

func EvaluateReviewQuality(pr *github.PullRequest, client clients.GithubRepoClient) error {
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
	for _, commentDetails := range filteredComments {

		comment := commentDetails.GetBody()
		diffHunk := commentDetails.GetDiffHunk()

		openaiClient := clients.NewOpenAiClient()
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
				fmt.Printf("Score:\t\t%d - Deleting\nComment:\t\t'%s'\n.", score, comment)
				err := client.DeleteComment(commentDetails)
				if err != nil {
					return err
				}
				continue
			} else if len(final) > 1 {
				// reply to the comment
				fmt.Printf("Score:\t\t%d - Keeping\nComment:\t\t'%s'\n", score, comment)
				err := client.ReplyToComment(pr, commentDetails, fmt.Sprintf("Thank you for your feedback! The AI gives your comment a score of %s", final))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
