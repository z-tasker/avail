package image

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
	"image/png"
	"os"
	"strings"
)

func MakeLogo(artStyles []string, logoPrompt string, tagline string, client *openai.Client, outputDir string) error {

	// Prepare the output directory
	fileInfo, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		os.MkdirAll(outputDir, os.ModePerm)
	} else {
		if !(fileInfo.IsDir()) {
			return errors.New(fmt.Sprintf("Error: outputDir '%s' exists and is not a directory", outputDir))
		}
	}

	generatedImageCount := 0

	for _, artStyle := range artStyles {
		log.Info().Msg(fmt.Sprintf("Generating set of '%s' images", artStyle))
		respBase64, err := client.CreateImage(
			context.Background(),
			openai.ImageRequest{
				Prompt:         fmt.Sprintf("%s %s %s", artStyle, logoPrompt, tagline),
				Size:           openai.CreateImageSize256x256,
				ResponseFormat: openai.CreateImageResponseFormatB64JSON,
				N:              3,
			},
		)

		if err != nil {
			return errors.Wrap(err, "Error during CreateImage")
		}

		for i, generatedImage := range respBase64.Data {

			imgBytes, err := base64.StdEncoding.DecodeString(generatedImage.B64JSON)
			if err != nil {
				return errors.Wrap(err, "Error during base64 decoding")
			}

			r := bytes.NewReader(imgBytes)
			imgData, err := png.Decode(r)
			if err != nil {
				return errors.Wrap(err, "Error during PNG decoding")
			}

			artStyleShort, _, _ := strings.Cut(artStyle, " ")

			file, err := os.Create(fmt.Sprintf("%s/logo-%s-%d.png", outputDir, artStyleShort, i))
			if err != nil {
				return errors.Wrap(err, "Error creating output file")
			}
			defer file.Close()

			if err := png.Encode(file, imgData); err != nil {
				return errors.Wrap(err, "Error during PNG encoding to file")
			}

			generatedImageCount++
		}
	}
	log.Info().Msg(fmt.Sprintf("%d candidate logos generated here: %s\n", generatedImageCount, outputDir))
	return nil
}
