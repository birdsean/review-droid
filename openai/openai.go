package openai

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

const (
	systemMessage = `
		You are an expert GitHub code reviewer. You are reviewing a pull request, and will be given snippets from the raw diff.
		If you see no problems, respond with "No comments".
		If you are certain there is a bug, prefix your comment with "Bug:".
		If the problem is a potential bug, do not comment.
		If the problem is a style issue, prefix your comment with "Style:".
		If the problem is a question, prefix your comment with "Question:".
		If the problem is a suggestion, prefix your comment with "Suggestion:".
		If the problem is a request for clarification, prefix your comment with "Clarification:".
		If a unit test of critical functionality is missing, prefix your comment with "Missing Test:".
		If a unit test could use some more test cases, prefix your comment with "Suggested Test Cases:".
		Copy the "+" or "-" into your comment prefix before the line number. 
		Only rarely on a line that starts with "-".
		An example response would look like this:
			[- Line 2] Bug: 'countPizzas' is being used elsewhere and still needs to be initialized
			[+ Line 42] Readability: consider saving this magic number to a variable
			[+ Line 43] Refactor Suggestion: This code is duplicated in 3 places. Consider refactoring into a function.
		Do not nitpick. Comments must be high quality and pithy.	
	`
	CODE_PREVIEW_SIZE = 4
)

type OpenAiClient struct {
	client *openai.Client
}

func (oac *OpenAiClient) Init() {
	openaiToken := os.Getenv("OPENAI_TOKEN")
	if openaiToken == "" {
		log.Fatalf("OPENAI_TOKEN environment variable must be set")
	}
	oac.client = openai.NewClient(openaiToken)
}

func (oac *OpenAiClient) GetCompletion(prompt string) (*string, error) {

	logMsg := "Evaluating code lines:\n"
	splitLines := strings.Split(prompt, "\n")
	for i, line := range splitLines {
		if i < 4 || i > len(splitLines)-4 {
			logMsg += line + "\n"
		} else if i == 4 {
			logMsg += "...\n"
		}
	}
	fmt.Print(logMsg)

	resp, err := oac.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemMessage,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	fmt.Printf(
		"PromptTokens:\t\t%d\nCompletionTokens:\t\t%d\nTotalTokens:\t\t%d\n",
		resp.Usage.PromptTokens,
		resp.Usage.CompletionTokens,
		resp.Usage.TotalTokens,
	)
	return &resp.Choices[0].Message.Content, nil
}
