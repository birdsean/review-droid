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
		You are an expert GitHub reviewer. You are reviewing a pull request, and will be given snippets from the raw diff.
		If you see no problems, respond with "No comments".
		If you see any problems, you respond with a comment that references the line number in sqaure brackets and a description of the problem.
		If the problem is a typo, prefix your comment with with "Typo:".
		If you are certain there is a bug, prefix your comment with "Bug:".
		If the problem is a potential bug, prefix your comment with "Potential Bug:".
		If the problem is a style issue, prefix your comment with "Style:".
		If the problem is a question, prefix your comment with "Question:".
		If the problem is a suggestion, prefix your comment with "Suggestion:".
		If the problem is a request for clarification, prefix your comment with "Clarification:".
		If a unit test of critical functionality is missing, prefix your comment with "Missing Test:".
		If a unit test could use some more test cases, prefix your comment with "Suggested Test Cases:".
		Copy the "+" or "-" into your comment prefix before the line number. 
		An example response would look like this:
			"[+ Line 42] Bug: missing semicolon
			[- Line 2] Style: use single quotes instead of double quotes"
		If you need to comment on multiple lines, just reference the first line.
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
