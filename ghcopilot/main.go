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
	// REQUIRED - The GitHub Token to authenticate with Copilot - usually a PAT with "Copilot Requests" with "Allow: Read Only: scope
	Token *dagger.Secret
	// OPTIONAL - The model to use for Copilot (e.g., claude-sonnet-4.5", "claude-sonnet-4", "gpt-5" defaults to the @github/copilot cli versions' default)
	// +optional
	Model string
	// OPTIONAL at constructiion - The prompt to send to Copilot
	Prompt string
}

func (c *Ghcopilot) NewGhcopilot(
	ctx context.Context,
	// The model to use for Copilot (e.g., claude-sonnet-4.5", "claude-sonnet-4", "gpt-5" defaults to the @github/copilot cli versions' default)
	// +optional
	model string,
	token *dagger.Secret,
) (*Ghcopilot, error ){

	if token == nil {
		return nil, fmt.Errorf("missing token secret: call ghcopilot with-token --token env:GITHUB_TOKEN (or your env var) before calling response")
	}
	return &Ghcopilot{
		Token: token,
		Model: model,
	}, nil
}

func (c *Ghcopilot) WithModel(
	ctx context.Context,
	model string,
) *Ghcopilot {
	return &Ghcopilot{
		Token:  c.Token,
		Model:  model,
		Prompt: c.Prompt,
	}
}

func (c *Ghcopilot) WithPrompt(
	ctx context.Context,
	// REQUIRED - The prompt to send to Copilot
	prompt string,
) *Ghcopilot {
	return &Ghcopilot{
		Token:  c.Token,
		Model:  c.Model,
		Prompt: prompt,
	}
}

func (c *Ghcopilot) Container(
	ctx context.Context,
) *dagger.Container {
	return dag.Container().
		From("node:alpine3.22").
		WithWorkdir("/workspace").
		WithExec([]string{"npm", "install", "-g", "@github/copilot"}).
		WithSecretVariable("GITHUB_TOKEN", c.Token)
}


// Returns a container with GitHub Copilot Installed
func (c *Ghcopilot) Response(
	ctx context.Context,
) (string, error) {
	container := c.Container(ctx)

	if c.Model != "" {
		container = container.WithExec([]string{"copilot", "--model", c.Model, "--prompt", c.Prompt})
	} else {
		container = container.WithExec([]string{"copilot", "--prompt", c.Prompt})
	}

	return container.Stdout(ctx)
}
