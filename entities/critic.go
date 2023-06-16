package entities

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/birdsean/review-droid/clients"
	"github.com/google/go-github/v53/github"
	"github.com/sashabaranov/go-openai"
)

type Critic struct {
	githubClient clients.GithubRepoClient
}

func NewCritic(client clients.GithubRepoClient) *Critic {
	return &Critic{client}
}

func (critic *Critic) EvaluateReviewQuality(prId int) error {
	// List review comments on a pull request
	comments, err := critic.githubClient.GetPrComments(prId)
	if err != nil {
		log.Fatalf("Failed to get comments: %v", err)
	}

	// filter out comments that have been resolved
	filteredComments := FilterReplies(comments)

	fmt.Printf("Evaluating %d comments. %d were filtered.\n", len(filteredComments), len(comments)-len(filteredComments))
	for _, comment := range filteredComments {
		err := critic.evaluateComment(comment, prId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (critic *Critic) evaluateComment(comment *github.PullRequestComment, prId int) error {
	commentBody := comment.GetBody()
	diffHunk := comment.GetDiffHunk()

	// read evaluate.txt file
	file, err := os.ReadFile("prompts/evaluate.v1.txt")
	if err != nil {
		log.Fatalf("Failed to read system message: %v", err)
	}
	evaluatePrompt := string(file)
	conversation, err := getCodeSummary(evaluatePrompt, diffHunk)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	score, err := getCommentQualityRating(conversation, commentBody)

	if err != nil {
		return err
	}

	if score < 3 {
		fmt.Printf("Score:\t\t%d - Deleting\nComment:\t\t'%s'\n.", score, commentBody)
		err := critic.githubClient.DeleteComment(comment)
		if err != nil {
			return err
		}
	}
	return nil
}

func getCodeSummary(evaluatePrompt, diffHunk string) ([]openai.ChatCompletionMessage, error) {
	openaiClient := clients.NewOpenAiClient()
	conversation := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: evaluatePrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Please summarize what this code is doing:\n%s", diffHunk),
		},
	}
	completion, err := openaiClient.RequestCompletion(conversation)
	if err != nil {
		return nil, err
	}
	conversation = append(conversation, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: completion,
	})
	return conversation, nil
}

func getCommentQualityRating(
	conversation []openai.ChatCompletionMessage,
	commentBody string,
) (int, error) {
	conversation = append(conversation, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: fmt.Sprintf("Please rate the quality of this code review comment on a scale of 1-5:\n%s\nOnly respond with the number.", commentBody),
	})
	openaiClient := clients.NewOpenAiClient()
	message, err := openaiClient.RequestCompletion(conversation)
	if err != nil {
		return -1, err
	}

	firstChar := message[0:1]
	if strings.Contains("12345", firstChar) {
		score, err := strconv.Atoi(firstChar)
		if err != nil {
			return -1, err
		}
		return score, nil
	} else {
		return -1, errors.New("No score detected when rating comment quality")
	}
}
