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

func MakeTestsUsage() string {
	return "Generates test files for a set of source files."
}

func MakeTestsFlags() []cli.Flag {
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

func MakeTests(ctx *cli.Context, config *config.Config, openAIClient *openai.Client) error {

	language, err := common.DetermineProgrammingLanguage(config, ctx.Args().Slice(), ctx.Int("MaxTokens"), openAIClient)
	if err != nil {
		return errors.Wrap(err, "Error while determining programming language")
	}
	log.Debug().Msg(fmt.Sprintf("Determined language as '%s'", language))
	codeStyleGuide, err := common.GetLanguageSpecificCodeStyleGuide(config, language)
	if err != nil {
		return errors.Wrap(err, "Error while constructing programming language specific code style guide")
	}
	packageName := ctx.String("PackageName")
	if len(packageName) > 0 {
		codeStyleGuide = fmt.Sprintf("%s. The package is called '%s'", codeStyleGuide, packageName)
	}

	tests, err := completion.ChatCompletionFromProjectFiles(
		config.MakeTestsPrompt(),
		ctx.Args().Slice(),
		openAIClient,
		codeStyleGuide,
		ctx.Int("MaxTokens"))
	if err != nil {
		return errors.Wrap(err, "Error generating tests")
	}
	fmt.Println(tests)

	return nil
}
