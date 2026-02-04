package cli

// ExitError wraps an error with an exit code for CLI flows.
type ExitError struct {
	Err  error
	Code int
}

// Error returns the underlying error message.
func (e ExitError) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

// Unwrap exposes the underlying error for errors.Is/As.
func (e ExitError) Unwrap() error {
	return e.Err
}

// ExitCode returns the process exit code for the error.
func (e ExitError) ExitCode() int {
	return e.Code
}
