// A generated module for Ghcopilot functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/ghcopilot/internal/dagger"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Ghcopilot struct {
	// REQUIRED - The GitHub PAT Token to authenticate with Copilot - must have permissions "Copilot Requests" with "Allow: Read Only" scope
	Token *dagger.Secret
	// OPTIONAL - The model to use for Copilot (e.g., claude-sonnet-4.5", "claude-sonnet-4", "gpt-5" defaults to the @github/copilot cli versions' default)
	Model string
	// OPTIONAL at constructiion - The prompt to send to Copilot
	Prompt string
	// OPTIONAL - The workspace directory to use as the context for Copilot (defaults to the root of the project)
	Workspace *dagger.Directory
}

type LLMResponse struct {
	Content    string
	TokenUsage LLMTokenUsage
}

type LLMTokenUsage struct {
	InputTokens       int64
	OutputTokens      int64
	CachedTokenReads  int64
	CachedTokenWrites int64
	TotalTokens       int64
}

func (c *Ghcopilot) NewGhcopilot(
	ctx context.Context,
	// The model to use for Copilot (e.g., claude-sonnet-4.5", "claude-sonnet-4", "gpt-5" defaults to the @github/copilot cli versions' default)
	// +optional
	model string,
	token *dagger.Secret,
	// +defaultPath="/"
	workspace *dagger.Directory,
) (*Ghcopilot, error) {

	if token == nil {
		return nil, fmt.Errorf("missing token secret: call ghcopilot with-token --token env:GITHUB_TOKEN (or your env var) before calling response")
	}

	return &Ghcopilot{
		Token:     token,
		Model:     model,
		Workspace: workspace,
	}, nil
}

func (c *Ghcopilot) WithToken(
	ctx context.Context,
	// REQUIRED - The GitHub PAT Token to authenticate with Copilot - must have permissions "Copilot Requests" with "Allow: Read Only" scope
	token *dagger.Secret,
) *Ghcopilot {

	c.Token = token
	return c
}

func (c *Ghcopilot) WithModel(
	ctx context.Context,
	model string,
) *Ghcopilot {

	c.Model = model
	return c
}

func (c *Ghcopilot) WithPrompt(
	ctx context.Context,
	// REQUIRED - The prompt to send to Copilot
	prompt string,
) *Ghcopilot {

	c.Prompt = prompt
	return c
}

// Returns a container with GitHub Copilot Installedfunc (c *Ghcopilot) Container(
func (c *Ghcopilot) Container(
	ctx context.Context,
) *dagger.Container {
	return dag.Container().
		From("node:24-bookworm-slim").
		WithExec([]string{"npm", "install", "-g", "@github/copilot"}).
		WithSecretVariable("GITHUB_TOKEN", c.Token).
		WithDirectory("/workspace", c.Workspace).
		WithWorkdir("/workspace")
}

// Runs Copilot with the given prompt
func (c *Ghcopilot) Response(
	ctx context.Context,
) (*LLMResponse, error) {

	container := c.Container(ctx)

	var content string
	var tokenUsage LLMTokenUsage

	if c.Model != "" {
		container = container.WithExec([]string{"copilot", "--model", c.Model, "--prompt", c.Prompt})
	} else {
		container = container.WithExec([]string{"copilot", "--prompt", c.Prompt})
	}

	content, err := container.Stdout(ctx)
	if err != nil {
		return nil, err
	}

	// Get the stderr output which contains usage information
	responseMetadata, err := container.Stderr(ctx)
	if err != nil {
		return nil, err
	}

	// Parse the token usage from the stderr output
	tokenUsage = parseLLMTokenUsage(responseMetadata)

	return &LLMResponse{
		Content:    content,
		TokenUsage: tokenUsage,
	}, nil
}

// parseLLMTokenUsage parses the stderr output from GitHub Copilot CLI to extract token usage information
func parseLLMTokenUsage(output string) LLMTokenUsage {
	var tokenUsage LLMTokenUsage

	// Parse the usage line that contains model-specific token information
	// Example: "claude-sonnet-4.5    7.5k input, 52 output, 3.6k cache read (Est. 1 Premium request)"

	// Look for the pattern: model name followed by input, output, cache read values
	re := regexp.MustCompile(`(\d+(?:\.\d+)?)(k?)\s+input,\s*(\d+(?:\.\d+)?)(k?)\s+output,\s*(\d+(?:\.\d+)?)(k?)\s+cache read,\s*(\d+(?:\.\d+)?)(k?)\s+cache write`)
	matches := re.FindStringSubmatch(output)

	if len(matches) >= 7 {
		// Parse input tokens
		if inputVal, err := strconv.ParseFloat(matches[1], 64); err == nil {
			if strings.ToLower(matches[2]) == "k" {
				inputVal *= 1000
			}
			tokenUsage.InputTokens = int64(inputVal)
		}

		// Parse output tokens
		if outputVal, err := strconv.ParseFloat(matches[3], 64); err == nil {
			if strings.ToLower(matches[4]) == "k" {
				outputVal *= 1000
			}
			tokenUsage.OutputTokens = int64(outputVal)
		}

		// Parse cache read tokens
		if cacheVal, err := strconv.ParseFloat(matches[5], 64); err == nil {
			if strings.ToLower(matches[6]) == "k" {
				cacheVal *= 1000
			}
			tokenUsage.CachedTokenReads = int64(cacheVal)
		}

		// Parse cache write tokens
		if cacheVal, err := strconv.ParseFloat(matches[7], 64); err == nil {
			if strings.ToLower(matches[8]) == "k" {
				cacheVal *= 1000
			}
			tokenUsage.CachedTokenWrites = int64(cacheVal)
		}

		tokenUsage.TotalTokens = tokenUsage.InputTokens + tokenUsage.OutputTokens
	}

	return tokenUsage
}
