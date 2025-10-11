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

type Copilot struct {
}

// Returns a container with GitHub Copilot Installed
func (m *Copilot) GoCode(
	ctx context.Context,
	model string,
	token *dagger.Secret,
	prompt string,

) (string, error) {

	container := dag.Container().
		From("node:alpine3.22").
		WithWorkdir("/workspace").
		WithExec([]string{"npm", "install", "-g", "@github/copilot"}).
		WithSecretVariable("GITHUB_TOKEN", token).
		WithExec([]string{"copilot", "--model", model, "--prompt", prompt})

	return container.Stdout(ctx)
}
