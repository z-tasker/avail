package cmd

import (
    "bufio"
	"fmt"
    "github.com/briandowns/spinner"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/z-tasker/avail/ai/completion"
	"github.com/z-tasker/avail/cmd/common"
	"github.com/z-tasker/avail/config"
	"github.com/z-tasker/avail/util"
	openai "github.com/sashabaranov/go-openai"
    "os"
    "strings"
    "time"
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
            Aliases: []string{"f"},
			Usage: "Optionally provide a file to include in the prompt",
		},
		&cli.StringFlag{
			Name:  "Language",
            Aliases: []string{"l"},
			Usage: "Optionally provide a language to add a corresponding style guide to the prompt",
		},
		&cli.BoolFlag{
			Name:  "Interactive",
            Aliases: []string{"i"},
			Usage: "Optionally enter a REPL where history is replayed on each prompt, mimicking chatGPT web interface",
		},
	}
}

func Prompt(ctx *cli.Context, config *config.Config, openAIClient *openai.Client) error {

    var promptBuilder strings.Builder

    // Get the bare prompt
	userPrompt := strings.Join(ctx.Args().Slice(), " ")
	if len(userPrompt) == 0 {
		return errors.New("Missing prompt")
	}
    promptBuilder.WriteString(userPrompt)

    // Get a language style guide if Language is set
    styleGuide, err := common.GetLanguageSpecificCodeStyleGuide(config, ctx.String("Language"))
    if err != nil {
        return errors.Wrap(err, "Error getting language specific style guide")
    }

    // Attach file contents to prompt if a file is provided
    includeFile := ctx.String("IncludeFile")
    if len(includeFile) > 0 {
        fileContent, err := util.FileContentsToPrompt("", []string{includeFile}, "")
        if err != nil {
            return errors.Wrap(err, fmt.Sprintf("Error adding file %s to prompt"))
        }
        promptBuilder.WriteString(fileContent)
    }
    
    log.Debug().
        Str("Language", ctx.String("Language")).
        Str("IncludeFile", ctx.String("IncludeFile")).
        Str("Prompt", userPrompt).
        Int("MaxTokens", ctx.Int("MaxTokens")).
        Msg("submitting prompt")
    // Submit the prompt
    messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: styleGuide,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: promptBuilder.String(),
		},
    }
    if ctx.Bool("Interactive") {
        colorGreen := "\033[32m"
        colorBlue := "\033[34m"
        colorCyan := "\033[36m"
        spin := spinner.New(spinner.CharSets[42], 100*time.Millisecond)
        spin.Prefix = fmt.Sprintf("%s 𝐴∀", string(colorCyan))
        spin.Color("cyan")
        reader := bufio.NewReader(os.Stdin)
        fmt.Println(string(colorBlue), "Conversation")
        fmt.Println(string(colorBlue), "------------------------")
        fmt.Printf("%s --> %s\n", string(colorGreen), promptBuilder.String())
        for {
            spin.Start()
	        resp, err := completion.ChatCompletion(messages, openAIClient, ctx.Int("MaxTokens"))
            spin.Stop()
	        if err != nil {
	        	return errors.Wrap(err, "Error retreiving ChatCompletion")
	        }
            fmt.Println(string(colorCyan), "𝐴∀:", resp)
            fmt.Println(string(colorBlue), "------------------------")
            messages = append(messages, openai.ChatCompletionMessage{
                Role: openai.ChatMessageRoleAssistant,
                Content: resp,
            })

            fmt.Println("")
            fmt.Print(string(colorGreen), " --> ")
            text, _ := reader.ReadString('\n')
            text = strings.Replace(text, "\n", "", -1)
            messages = append(messages, openai.ChatCompletionMessage{
                Role: openai.ChatMessageRoleUser,
                Content: text,
            })
        }
    } else {
	    resp, err := completion.ChatCompletion(messages, openAIClient, ctx.Int("MaxTokens"))
	    if err != nil {
	    	return errors.Wrap(err, "Error retreiving ChatCompletion")
	    }
	    fmt.Println(resp)
    }

	return nil
}
