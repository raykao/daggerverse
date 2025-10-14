// A generated module for Copilot functions
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
	"dagger/copilot/internal/dagger"
)

type GHCPClient struct {
	// REQUIRED - The GitHub Token to authenticate with Copilot - usually a PAT with "Copilot Requests" with "Allow: Read Only: scope
	token *dagger.Secret
	// OPTIONAL - The model to use for Copilot (e.g., claude-sonnet-4.5", "claude-sonnet-4", "gpt-5" defaults to the @github/copilot cli versions' default)
	// +optional
	Model string
	// OPTIONAL at constructiion - The prompt to send to Copilot
	Prompt string
}

func NewGHCPClient(
	ctx context.Context,
	// The model to use for Copilot (e.g., claude-sonnet-4.5", "claude-sonnet-4", "gpt-5" defaults to the @github/copilot cli versions' default)
	// +optional
	model string,
	token *dagger.Secret,
) *GHCPClient {
	return &GHCPClient{
		token: token,
		Model: model,
	}
}

func (c *GHCPClient) WithPrompt(
	// REQUIRED - The prompt to send to Copilot
	prompt string,
) *GHCPClient {
	return &GHCPClient{
		token:  c.token,
		Model:  c.Model,
		Prompt: prompt,
	}
}

// Returns a container with GitHub Copilot Installed
func (c *GHCPClient) Response(
	ctx context.Context,
) (string, error) {
	container := dag.Container().
		From("node:alpine3.22").
		WithWorkdir("/workspace").
		WithExec([]string{"npm", "install", "-g", "@github/copilot"}).
		WithSecretVariable("GITHUB_TOKEN", c.token)

	if c.model != "" {
		container = container.WithExec([]string{"copilot", "--model", c.model, "--prompt", c.prompt})
	} else {
		container = container.WithExec([]string{"copilot", "--prompt", c.prompt})
	}

	return container.Stdout(ctx)
}
