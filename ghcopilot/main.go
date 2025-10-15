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
	"fmt"
	"context"
	"dagger/ghcopilot/internal/dagger"
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
	Info string
}

func (c *Ghcopilot) NewGhcopilot(
	ctx context.Context,
	// The model to use for Copilot (e.g., claude-sonnet-4.5", "claude-sonnet-4", "gpt-5" defaults to the @github/copilot cli versions' default)
	// +optional
	model string,
	token *dagger.Secret,
	// +defaultPath="/"
	workspace *dagger.Directory,
) (*Ghcopilot, error ){

	if token == nil {
		return nil, fmt.Errorf("missing token secret: call ghcopilot with-token --token env:GITHUB_TOKEN (or your env var) before calling response")
	}

	return &Ghcopilot{
		Token: token,
		Model: model,
		Workspace: workspace,
	}, nil
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

	// Note: The GitHub Copilot CLI does not currently provide token usage details in its output.
	// If it did, you would parse that information here and populate the TokenUsage struct accordingly.
	info, err := container.Stderr(ctx)
	if err != nil {
		return nil, err
	}

	tokenUsage = LLMTokenUsage{
		Info: info,
	}

	return &LLMResponse{
		Content: content,
		TokenUsage: tokenUsage,
	}, nil
}