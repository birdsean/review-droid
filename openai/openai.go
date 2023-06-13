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
		Lines that start with "-" are lines that have been removed. Lines that start with "+" are lines that have been added.
		If you see no problems, respond with "No comments".
		If you are certain there is a bug, prefix your comment with "Bug:".
		If the problem is a potential bug, do not comment.
		If the problem is a style issue, do not comment.
		If the problem is a question, prefix your comment with "Question:".
		If the problem is a suggestion, prefix your comment with "Suggestion:".
		If the problem is a request for clarification, prefix your comment with "Clarification:".
		If a unit test of critical functionality is missing, prefix your comment with "Missing Test:".
		If a unit test could use some more test cases, prefix your comment with "Suggested Test Cases:".
		If a method or class is too big, prefix your comment with "Refactor Suggestion:".
		If you see lots of duplicated code, prefix your comment with "Refactor Suggestion:".
		Copy the "+" or "-" into your comment prefix before the line number. 
		Only rarely comment on a line that starts with "-".
		Do not comment on imports.
		Do not nitpick. Comments must be high quality and pithy.
		You can comment on multiple lines.	
		An example response would look like this:
			[- Line 2] Bug: 'countPizzas' is being used elsewhere and still needs to be initialized
			[+ Line 42] Readability: consider saving this magic number to a variable
			[+ Line 43] Refactor Suggestion: This code is duplicated in 3 places. Consider refactoring into a function.
	`
	CODE_PREVIEW_SIZE = 4
	followUpMessage   = `Please make correct any line references you may have gone wrong.
	If you answered "No comments", respond with "No comments" again.
	Remove low quality comments, if you remove all comments, respond with "No comments". 
	Rewrite your comments if they need it. 
	REMOVE ALL COMMENTS that have to do with IMPORT STATEMENTS.
	Maintain the same format as your first response.
	If no changes are needed from your last response, respond with "No changes".`
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

func printTokenUsage(response openai.ChatCompletionResponse, countInputMessages int) {
	content := response.Choices[0].Message.Content
	if len(content) > 100 {
		content = content[:100]
	}
	fmt.Printf(
		"CountInputMessages:\t%d\nTotalTokens:\t\t%d\nPreview:\t\t%s\n******************************\n",
		countInputMessages,
		response.Usage.TotalTokens,
		content,
	)
}

func (oac *OpenAiClient) RequestCompletion(messages []openai.ChatCompletionMessage) (string, error) {
	resp, err := oac.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	if err != nil {
		return "", err
	}

	printTokenUsage(resp, len(messages))
	return resp.Choices[0].Message.Content, nil
}

func (oac *OpenAiClient) GetCompletion(prompt string, debug bool) (*string, error) {

	if debug {
		logMsg := "Evaluating code lines:\n"
		splitLines := strings.Split(prompt, "\n")
		for i, line := range splitLines {
			if i < CODE_PREVIEW_SIZE || i > len(splitLines)-CODE_PREVIEW_SIZE {
				logMsg += line + "\n"
			} else if i == CODE_PREVIEW_SIZE {
				logMsg += "...\n"
			}
		}
		fmt.Print(logMsg)
	}

	firstDraft, err := oac.RequestCompletion([]openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemMessage,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	})
	if err != nil {
		return nil, err
	}

	secondDraft, err := oac.RequestCompletion([]openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemMessage,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
		{
			Role:    openai.ChatMessageRoleAssistant,
			Content: firstDraft,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: followUpMessage,
		},
	})
	if err != nil {
		return nil, err
	}

	if strings.Contains(secondDraft, "No changes") {
		return &firstDraft, nil
	}
	return &secondDraft, nil
}
