package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"
	"github.com/z-tasker/avail/ai/completion"
	"github.com/z-tasker/avail/ai/image"
	"github.com/z-tasker/avail/config"
)

func MakeLogoUsage() string {
	return "Generates a set of candidate logos for the project implemented by a set of source files."
}

func MakeLogoFlags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:  "MaxTokens",
			Value: 5000,
		},
		&cli.StringFlag{
			Name:     "OutputDir",
			Usage:    "Directory where logos will be saved",
			Required: true,
		},
	}
}

func MakeLogo(ctx *cli.Context, config *config.Config, openAIClient *openai.Client) error {

	// First get a tagline, then use it to generate logos
	tagline, err := completion.ChatCompletionFromProjectFiles(
		config.MakeTaglinePrompt(),
		ctx.Args().Slice(),
		openAIClient,
		"",
		ctx.Int("MaxTokens"))
	log.Info().Msg(fmt.Sprintf("Generated tagline for logo generation: %s", tagline))

	if err != nil {
		return errors.Wrap(err, "Error generating tagline")
	}

	err = image.MakeLogo(
		config.MakeLogoArtStyles(),
		config.MakeLogoPrompt(),
		tagline,
		openAIClient,
		ctx.String("OutputDir"))

	return err
}
