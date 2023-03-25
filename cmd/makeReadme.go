package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"
	"github.com/z-tasker/avail/ai/completion"
	"github.com/z-tasker/avail/cmd/common"
	"github.com/z-tasker/avail/config"
)

func MakeReadmeUsage() string {
	return "Generates a README file for a set of source files. The parent directory of the first file provided will be used as the 'name' of the project."
}

func MakeReadmeFlags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:  "MaxTokens",
			Value: 6000,
		},
		&cli.StringFlag{
			Name:  "PackageName",
			Usage: "Optionally provide a package name to help the generated tests be more complete",
		},
	}
}

func MakeReadme(ctx *cli.Context, config *config.Config, openAIClient *openai.Client) error {

	language, err := common.DetermineProgrammingLanguage(config, ctx.Args().Slice(), ctx.Int("MaxTokens"), openAIClient)
	if err != nil {
		return errors.Wrap(err, "Error while determining programming language")
	}
	log.Debug().Msg(fmt.Sprintf("Determined language as '%s'", language))
	readmeStyleGuide, err := common.GetLanguageSpecificReadmeStyleGuide(config, language)
	if err != nil {
		return errors.Wrap(err, "Error while constructing programming language specific README style guide")
	}
	packageName := ctx.String("PackageName")
	if len(packageName) > 0 {
		readmeStyleGuide = fmt.Sprintf("%s. The project is called '%s'", readmeStyleGuide, packageName)
	}

	readme, err := completion.ChatCompletionFromProjectFiles(
		config.MakeReadmePrompt(),
		ctx.Args().Slice(),
		openAIClient,
		readmeStyleGuide,
		ctx.Int("MaxTokens"))
	if err != nil {
		return errors.Wrap(err, "Error while generating Readme")
	}
	fmt.Println(readme)

	return nil
}
