package completion

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
	"github.com/z-tasker/avail/util"
	"path"
	"path/filepath"
)

func SimpleCompletion(prompt string, openAIClient *openai.Client, maxTokens int) (string, error) {
	resp, err := openAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT4,
			MaxTokens: maxTokens,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	completion := resp.Choices[0].Message.Content

	log.Debug().
		Str("Prompt", prompt).
		Str("Response", completion).
		Int("MaxTokens", maxTokens).
		Msg("")

	return completion, nil
}

func ChatCompletionFromProjectFiles(promptPrefix string, sourceFiles []string, openAIClient *openai.Client, styleGuide string, maxTokens int) (string, error) {

	firstPath, err := filepath.Abs(path.Dir(sourceFiles[0]))
	if err != nil {
		return "", err
	}
	firstPathParentDir := path.Base(firstPath)

	if len(sourceFiles) < 1 {
		return "", errors.New("Missing sourceFiles arguments")
	}

	prompt, err := util.FileContentsToPrompt(promptPrefix, sourceFiles, fmt.Sprintf("This project is called '%s'", firstPathParentDir))
	if err != nil {
		return "", errors.Wrap(err, "Error extracting sourceFiles contents to prompt for ChatCompletionFromProjectFiles")
	}

	resp, err := openAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT4,
			MaxTokens: maxTokens,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: styleGuide,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", errors.Wrap(err, "Error during CreateChatCompletion for ChatCompletionFromProjectFiles")
	}

	completion := resp.Choices[0].Message.Content

	log.Debug().
		Str("StyleGuide", styleGuide).
		Str("Prompt", prompt).
		Str("Response", completion).
		Int("MaxTokens", maxTokens).
		Msg("")

	return completion, nil
}
