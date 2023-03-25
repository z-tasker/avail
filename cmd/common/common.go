package common

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
	"github.com/z-tasker/avail/ai/completion"
	"github.com/z-tasker/avail/config"
	"strings"
)

func DetermineProgrammingLanguage(config *config.Config, sourceFiles []string, maxTokens int, openAIClient *openai.Client) (string, error) {

	language, err := completion.ChatCompletionFromProjectFiles(
		config.DetermineLanguagePrompt(),
		sourceFiles,
		openAIClient,
		"",
		maxTokens)

	if err != nil {
		return "", errors.Wrap(err, "Error while determining programming language from source files")
	}

	log.Debug().Msg(fmt.Sprintf("Determined language as '%s'", language))

	return language, nil

}

func GetLanguageSpecificReadmeStyleGuide(config *config.Config, language string) (string, error) {

	var readmeStyleGuide string

	switch strings.ToLower(language) {
	case "python":
		readmeStyleGuide = fmt.Sprintf("%s %s", config.ReadmeStyleGuideShared(), config.ReadmeStyleGuidePython())
	case "go":
		readmeStyleGuide = fmt.Sprintf("%s %s", config.ReadmeStyleGuideShared(), config.ReadmeStyleGuideGo())
	case "typescript":
		readmeStyleGuide = fmt.Sprintf("%s %s", config.ReadmeStyleGuideShared(), config.ReadmeStyleGuideTypescript())
	default:
		log.Warn().Msg(fmt.Sprintf("No style guide available for %s", language))
		readmeStyleGuide = config.ReadmeStyleGuideShared()
	}

	return readmeStyleGuide, nil

}

func GetLanguageSpecificCodeStyleGuide(config *config.Config, language string) (string, error) {

	var codeStyleGuide string

	switch strings.ToLower(language) {
	case "python":
		codeStyleGuide = fmt.Sprintf("%s %s", config.CodeStyleGuideShared(), config.CodeStyleGuidePython())
	case "go":
		codeStyleGuide = fmt.Sprintf("%s %s", config.CodeStyleGuideShared(), config.CodeStyleGuideGo())
	case "typescript":
		codeStyleGuide = fmt.Sprintf("%s %s", config.CodeStyleGuideShared(), config.CodeStyleGuideTypescript())
	default:
		log.Warn().Msg(fmt.Sprintf("No style guide available for %s", language))
		codeStyleGuide = config.CodeStyleGuideShared()
	}

	return codeStyleGuide, nil

}
