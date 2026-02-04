package workflow

// ConfigError wraps configuration errors for workflow use.
type ConfigError struct {
	Err error
}

// Error returns the underlying error message.
func (e ConfigError) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

// Unwrap exposes the underlying error for errors.Is/As.
func (e ConfigError) Unwrap() error {
	return e.Err
}

// ProviderError wraps provider errors for workflow use.
type ProviderError struct {
	Err error
}

// Error returns the underlying error message.
func (e ProviderError) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

// Unwrap exposes the underlying error for errors.Is/As.
func (e ProviderError) Unwrap() error {
	return e.Err
}
