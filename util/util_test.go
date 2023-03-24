package util_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/z-tasker/avail/util"
	"io/ioutil"
	"os"
	"testing"
)

func TestFileContentsToPrompt(t *testing.T) {
	err := ioutil.WriteFile("testFile1.txt", []byte("Hello"), 0644)
	assert.NoError(t, err)
	defer func() { _ = os.Remove("testFile1.txt") }()
	err = ioutil.WriteFile("testFile2.txt", []byte("World"), 0644)
	assert.NoError(t, err)
	defer func() { _ = os.Remove("testFile2.txt") }()

	t.Run("basic test", func(t *testing.T) {
		result, err := util.FileContentsToPrompt("prefix_", []string{"testFile1.txt", "testFile2.txt"}, "_suffix")
		assert.NoError(t, err)
		assert.Equal(t, "prefix_\ntestFile1.txt: ```Hello```\ntestFile2.txt: ```World```\n_suffix", result)
	})

	t.Run("error test", func(t *testing.T) {
		_, err := util.FileContentsToPrompt("error_", []string{"nonExistantFile.txt"}, "_suffix")
		assert.Error(t, err)

	})
}
