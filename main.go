package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"
	"github.com/z-tasker/avail/ai/completion"
	"github.com/z-tasker/avail/ai/image"
	"github.com/z-tasker/avail/config"
	"github.com/z-tasker/avail/logger"
	"os"
	"strings"
)

func main() {
	config, err := config.New()
	if err != nil {
		panic(errors.Wrap(err, "config.New failed"))
	}

	logger.Init(config)

	openAIClient := openai.NewClient(config.OpenAIAPIKey())

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "prompt",
				Aliases: []string{"p"},
				Usage:   "Sends a prompt to the OpenAI API",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "MaxTokens",
						Value: 1000,
					},
				},
				Action: func(cCtx *cli.Context) error {
					userPrompt := cCtx.Args().First()
					if len(userPrompt) == 0 {
						return errors.New("Missing prompt")
					}

					otherArgs := cCtx.Args().Tail()
					if len(otherArgs) > 0 {
						return errors.New("Unkown arguments, provide exactly 1 prompt")
					}

					resp, err := completion.SimpleCompletion(userPrompt, openAIClient, cCtx.Int("MaxTokens"))
					if err != nil {
						return errors.Wrap(err, "Error submitting prompt")
					}

					fmt.Println(resp)
					return nil
				},
			},
			{
				Name:  "MakeReadme",
				Usage: "Generates a README file for a set of source files. The parent directory of the first file provided will be used as the 'name' of the project.",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "MaxTokens",
						Value: 6000,
					},
				},
				Action: func(cCtx *cli.Context) error {

					readme, err := completion.ChatCompletionFromProjectFiles(
						config.MakeReadmePrompt(),
						cCtx.Args().Slice(),
						openAIClient,
						config.StyleGuideReadme(),
						cCtx.Int("MaxTokens"))

					if err != nil {
						return errors.Wrap(err, "Error generating README")
					}

					fmt.Println(readme)
					return nil
				},
			},
			{
				Name:  "MakeTests",
				Usage: "Generates test files for a set of source files.",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "MaxTokens",
						Value: 6000,
					},
					&cli.StringFlag{
						Name:  "PackageName",
						Usage: "Provide a package name to help the generated tests be more complete",
					},
				},
				Action: func(cCtx *cli.Context) error {

					language, err := completion.ChatCompletionFromProjectFiles(
						config.DetermineLanguagePrompt(),
						cCtx.Args().Slice(),
						openAIClient,
						"",
						cCtx.Int("MaxTokens"))

					log.Debug().Msg(fmt.Sprintf("Determined language as '%s'", language))

					var styleGuide string

					switch strings.ToLower(language) {
					case "python":
						styleGuide = config.StyleGuidePython()
					case "go":
						styleGuide = config.StyleGuideGo()
					case "typescript":
						styleGuide = config.StyleGuideTypescript()
					default:
						log.Warn().Msg(fmt.Sprintf("No style guide available for %s", language))
						styleGuide = ""
					}

					packageName := cCtx.String("PackageName")

					if len(packageName) > 0 {
						styleGuide = fmt.Sprintf("%s. The package is called '%s'", styleGuide, packageName)
					}

					tests, err := completion.ChatCompletionFromProjectFiles(
						config.MakeTestsPrompt(),
						cCtx.Args().Slice(),
						openAIClient,
						styleGuide,
						cCtx.Int("MaxTokens"))

					if err != nil {
						return errors.Wrap(err, "Error generating tests")
					}

					fmt.Println(tests)
					return nil
				},
			},
			{
				Name:  "MakeTagline",
				Usage: "Generates a tagline for the project implemented by a set of source files. The parent directory of the first file provided will be used as the 'name' of the project.",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "MaxTokens",
						Value: 5000,
					},
				},
				Action: func(cCtx *cli.Context) error {

					tagline, err := completion.ChatCompletionFromProjectFiles(
						config.MakeTaglinePrompt(),
						cCtx.Args().Slice(),
						openAIClient,
						config.StyleGuideReadme(),
						cCtx.Int("MaxTokens"))

					if err != nil {
						return errors.Wrap(err, "Error generating tagline")
					}

					fmt.Println(tagline)
					return nil
				},
			},
			{
				Name:  "MakeLogo",
				Usage: "Generates a set of candidate logos for the project implemented by a set of source files.",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "MaxTokens",
						Value: 5000,
					},
					&cli.StringFlag{
						Name:     "OutputDir",
						Usage:    "Directory where logos will be saved",
						Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {

					// First get a tagline, then use it to generate logos
					tagline, err := completion.ChatCompletionFromProjectFiles(
						config.MakeTaglinePrompt(),
						cCtx.Args().Slice(),
						openAIClient,
						config.StyleGuideReadme(),
						cCtx.Int("MaxTokens"))
					log.Info().Msg(fmt.Sprintf("Generated tagline for logo generation: %s", tagline))

					if err != nil {
						return errors.Wrap(err, "Error generating tagline")
					}

					err = image.MakeLogo(
						config.MakeLogoArtStyles(),
						config.MakeLogoPrompt(),
						tagline,
						openAIClient,
						cCtx.String("OutputDir"))

					return err
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Msg(fmt.Sprintf("%v", err))
	}
}
