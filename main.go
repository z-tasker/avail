package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"
	"github.com/z-tasker/avail/config"
	"github.com/z-tasker/avail/cmd"
	"github.com/z-tasker/avail/logger"
	"os"
)

func main() {
	config, err := config.New()
	if err != nil {
		panic(errors.Wrap(err, "config.New failed"))
	}

	logger.Init(config)

	openAIClient := openai.NewClient(config.OpenAIAPIKey())

	app := &cli.App{
        Name: "avail",
        Usage: "A Veritable AI Lackey to help streamline your development process.",
		Commands: []*cli.Command{
			{
				Name:    "Prompt",
				Aliases: []string{"p"},
				Usage:   cmd.PromptUsage(),
				Flags:   cmd.PromptFlags(),
				Action: func(ctx *cli.Context) error {
                    return cmd.Prompt(ctx, config, openAIClient)
				},
			},
			{
				Name:  "MakeReadme",
				Usage: cmd.MakeReadmeUsage(),
				Flags: cmd.MakeReadmeFlags(),
				Action: func(ctx *cli.Context) error {
                    return cmd.MakeReadme(ctx, config, openAIClient)
				},
			},
			{
				Name:  "MakeTests",
				Usage: cmd.MakeTestsUsage(),
				Flags: cmd.MakeTestsFlags(),
				Action: func(ctx *cli.Context) error {
                    return cmd.MakeTests(ctx, config, openAIClient)
				},
			},
			{
				Name:  "MakeTagline",
				Usage: cmd.MakeTaglineUsage(),
				Flags: cmd.MakeTaglineFlags(),
				Action: func(ctx *cli.Context) error {
                    return cmd.MakeTagline(ctx, config, openAIClient)
				},
			},
			{
				Name:  "MakeLogo",
				Usage: cmd.MakeLogoUsage(),
				Flags: cmd.MakeLogoFlags(),
				Action: func(ctx *cli.Context) error {
                    return cmd.MakeLogo(ctx, config, openAIClient)

				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Msg(fmt.Sprintf("%v", err))
	}
}
