package workflow

type ConfigError struct {
	Err error
}

func (e ConfigError) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

func (e ConfigError) Unwrap() error {
	return e.Err
}

type ProviderError struct {
	Err error
}

func (e ProviderError) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

func (e ProviderError) Unwrap() error {
	return e.Err
}
