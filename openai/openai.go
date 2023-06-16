package openai

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

const CODE_PREVIEW_SIZE = 4

func readMessageFromFile(filename string) string {
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read system message: %v", err)
	}
	return string(file)
}

// read system message from file
var systemMessage = readMessageFromFile("prompts/system.v1.txt")
var followUpMessage = readMessageFromFile("prompts/follow-up.v1.txt")

type OpenAiClient struct {
	client *openai.Client
}

// constructor
func NewOpenAiClient() *OpenAiClient {
	oac := &OpenAiClient{}
	oac.Init()
	return oac
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

func printCompletionInfo(prompt string) {
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

func (oac *OpenAiClient) GetCompletion(prompt string, debug bool) (*string, error) {

	if debug {
		printCompletionInfo(prompt)
	}

	conversation := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemMessage,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	firstDraft, err := oac.RequestCompletion(conversation)
	if err != nil {
		return nil, err
	}

	conversation = append(conversation, []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleAssistant,
			Content: firstDraft,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: followUpMessage,
		},
	}...)

	secondDraft, err := oac.RequestCompletion(conversation)
	if err != nil {
		return nil, err
	}

	if strings.Contains(secondDraft, "No changes") {
		return &firstDraft, nil
	}
	return &secondDraft, nil
}
