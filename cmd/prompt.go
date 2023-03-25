package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"
	"github.com/z-tasker/avail/ai/completion"
)

func PromptUsage() string {
	return "Sends a prompt to the OpenAI API"
}

func PromptFlags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:  "MaxTokens",
			Value: 1000,
		},
		&cli.StringFlag{
			Name:  "IncludeFile",
			Usage: "Optionally provide a file to include in the prompt",
		},
	}
}

func Prompt(ctx *cli.Context, openAIClient *openai.Client) error {

	userPrompt := ctx.Args().First()
	if len(userPrompt) == 0 {
		return errors.New("Missing prompt")
	}
	otherArgs := ctx.Args().Tail()
	if len(otherArgs) > 0 {
		log.Warn().Msg("Unkown arguments, provide exactly 1 prompt")
	}

	resp, err := completion.SimpleCompletion(userPrompt, openAIClient, ctx.Int("MaxTokens"))
	if err != nil {
		return errors.Wrap(err, "Error retreiving SimpleCompletion")
	}
	fmt.Println(resp)

	return nil
}
