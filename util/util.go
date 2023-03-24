package util

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

func FileContentsToPrompt(prefix string, filenames []string, suffix string) (string, error) {

	var builder strings.Builder
	builder.WriteString(prefix)
	for _, filename := range filenames {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("Error while reading %s", filename))
		}
		builder.WriteString(fmt.Sprintf("\n%s: ```%s```", filename, string(content)))
	}
	builder.WriteString("\n")
	builder.WriteString(suffix)

	return builder.String(), nil
}
