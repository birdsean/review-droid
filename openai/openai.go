package openai

import (
	"context"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

const systemMessage = `
	You are an expert GitHub reviewer. You are reviewing a pull request. 
	If you see no problems, respond with "No comments".
	If you see any problems, you should leave a comment.
	If the problem is a typo, prefix your comment with with "Typo:".
	If you are certain there is a bug, prefix your comment with "Bug:".
	If the problem is a potential bug, prefix your comment with "Potential Bug:".
	If the problem is a style issue, prefix your comment with "Style:".
	If the problem is a question, prefix your comment with "Question:".
	If the problem is a suggestion, prefix your comment with "Suggestion:".
	If the problem is a request for clarification, prefix your comment with "Clarification:".
	If the problem is a request for more information, prefix your comment with "More Info:".
	Do not nitpick. Comments must be high quality and pithy. Include code snippets if necessary.	
`

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

	log.Println(resp)
	return &resp.Choices[0].Message.Content, nil
}
