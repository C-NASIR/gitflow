package release

import (
	"errors"

	"gitflow/internal/cli"
	"gitflow/internal/workflow"
)

const (
	exitCodeComputation = 1
	exitCodeConfig      = 2
	exitCodeProvider    = 3
)

func releaseExitError(err error) error {
	if err == nil {
		return nil
	}
	var configErr workflow.ConfigError
	if errors.As(err, &configErr) {
		return cli.ExitError{Err: err, Code: exitCodeConfig}
	}
	var providerErr workflow.ProviderError
	if errors.As(err, &providerErr) {
		return cli.ExitError{Err: err, Code: exitCodeProvider}
	}
	return cli.ExitError{Err: err, Code: exitCodeComputation}
}
