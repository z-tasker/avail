package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	openai "github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"
	"github.com/z-tasker/avail/ai/completion"
	"github.com/z-tasker/avail/config"
)

func MakeTaglineUsage() string {
	return "Generates a tagline for the project implemented by a set of source files. The parent directory of the first file provided will be used as the 'name' of the project."
}

func MakeTaglineFlags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:  "MaxTokens",
			Value: 5000,
		},
	}
}

func MakeTagline(ctx *cli.Context, config *config.Config, openAIClient *openai.Client) error {

	tagline, err := completion.ChatCompletionFromProjectFiles(
		config.MakeTaglinePrompt(),
		ctx.Args().Slice(),
		openAIClient,
		"",
		ctx.Int("MaxTokens"))
	if err != nil {
		return errors.Wrap(err, "Error generating tagline")
	}
	fmt.Println(tagline)

	return nil
}
