package config

import (
	"github.com/Psiphon-Inc/configloader-go"
	"github.com/Psiphon-Inc/configloader-go/toml"
	"github.com/pkg/errors"
	"os/user"
	"path"
)

type availConfig struct {
	Log struct {
		Level  string
		Format string
	}

	StyleGuides struct {
		Readme     string
		Python     string
		Go         string
		Typescript string
	}

	DetermineLanguage struct {
		Prompt string
	}

	MakeTagline struct {
		Prompt string
	}

	MakeTests struct {
		Prompt string
	}

	MakeReadme struct {
		Prompt string
	}

	MakeLogo struct {
		ArtStyles []string
		Prompt    string
	}

	OpenAIAPIKey string
}

type Config struct {
	data     availConfig
	metadata configloader.Metadata
}

func New() (*Config, error) {
	var conf Config

	// get user homedir for search path
	currentUser, err := user.Current()
	if err != nil {
		return nil, errors.Wrap(err, "Error determining current user for home directory")
	}
	homeDir := currentUser.HomeDir

	//
	// load config
	//

	fileLocations := []configloader.FileLocation{
		{
			Filename:    "avail.toml",
			SearchPaths: []string{".", path.Join(homeDir, ".config"), "/etc/avail"},
		},
	}

	configReaders, configClosers, configReaderNames, err := configloader.FindFiles(fileLocations...)

	if err != nil {
		return nil, errors.Wrap(err, "configloader.FindFiles failed for config files")
	}

	defer func() {
		for _, r := range configClosers {
			r.Close()
		}
	}()

	defaults := []configloader.Default{
		{
			Key: configloader.Key{"Log", "Level"},
			Val: "info",
		},
		{
			Key: configloader.Key{"Log", "Format"},
			Val: "console",
		},
	}

	envOverrides := []configloader.EnvOverride{
		{
			EnvVar: "OPENAI_API_KEY",
			Key:    configloader.Key{"OpenAIAPIKey"},
		},
	}

	conf.metadata, err = configloader.Load(
		toml.Codec,
		configReaders, configReaderNames,
		defaults,
		envOverrides,
		&conf.data)
	if err != nil {
		return nil, errors.Wrap(err, "configloader.Load failed for config")
	}

	return &conf, nil
}

func (c *Config) OpenAIAPIKey() string {
	return c.data.OpenAIAPIKey
}

func (c *Config) StyleGuideReadme() string {
	return c.data.StyleGuides.Readme
}

func (c *Config) StyleGuidePython() string {
	return c.data.StyleGuides.Python
}

func (c *Config) StyleGuideGo() string {
	return c.data.StyleGuides.Go
}

func (c *Config) StyleGuideTypescript() string {
	return c.data.StyleGuides.Typescript
}

func (c *Config) DetermineLanguagePrompt() string {
	return c.data.DetermineLanguage.Prompt
}

func (c *Config) MakeReadmePrompt() string {
	return c.data.MakeReadme.Prompt
}

func (c *Config) MakeTaglinePrompt() string {
	return c.data.MakeTagline.Prompt
}

func (c *Config) MakeTestsPrompt() string {
	return c.data.MakeTests.Prompt
}

func (c *Config) MakeLogoArtStyles() []string {
	return c.data.MakeLogo.ArtStyles
}

func (c *Config) MakeLogoPrompt() string {
	return c.data.MakeLogo.Prompt
}

func (c *Config) LogLevel() string {
	return c.data.Log.Level
}

func (c *Config) LogFormat() string {
	return c.data.Log.Format
}
