# Avail

A Veritable AI Lackey (avail) is a command-line tool that automates various tasks for your projects such as generating README files, making tests, creating taglines, and producing logos. It leverages the OpenAI API to perform these tasks intelligently.

In most cases, the content the current language models produces is a suitable starting point, but often need to be edited and massaged into a working/higher quality state.

As the name implies, this tool is your lackey, you are still in charge! (for now)

![Logo](logo.png?raw=true)

## Requirements

- To use Avail, you need an [OpenAI API key](https://beta.openai.com/signup/).
- Go (v1.11+) for compiling and running Avail

## Installation

1. Clone the Avail repository:

```
git clone https://github.com/z-tasker/avail.git
```

2. Change to the Avail directory and build the project:

```
cd avail
go build
```

## Usage

Before running Avail, ensure that you have set the OpenAI API key as an environment variable:

```
export OPENAI_API_KEY=yourapikey
```

### Available Commands

1. `prompt`: Sends a prompt to the OpenAI API

```
avail prompt "Hello, how are you?"
```

2. `MakeReadme`: Generates a README file for a set of source files, for instance this README was initialized with:

```
avail MakeReadme main.go
```

3. `MakeTests`: Generates test files for a set of source files, for instance some of this repo's tests were initalized with:

```
avail MakeTests --PackageName github.com/z-tasker/avail util/util.go
```

4. `MakeTagline`: Generates a tagline for the project

```
avail MakeTagline main.go
```

5. `MakeLogo`: Generates a set of candidate logos for the project, you guessed it, the logo in this readme was generated with:

```
avail MakeLogo --OutputDir logos main.go
```

Note: Replace `main.go` with the path to your project's source files.

## Customization

You can customize the prompts and style guides for each command by modifying the `avail.toml` config file.
