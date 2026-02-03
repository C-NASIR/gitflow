package ui

import (
	"strings"

	"github.com/charmbracelet/huh"
)

type CommitPromptInput struct {
	Conventional bool

	Type     string
	Scope    string
	Summary  string
	Body     string
	Breaking bool

	All bool
}

func PromptCommit(types []string, scopes []string, requireScope bool, defaults CommitPromptInput) (CommitPromptInput, error) {
	in := defaults

	var typeOptions []huh.Option[string]
	for _, t := range types {
		typeOptions = append(typeOptions, huh.NewOption(t, t))
	}

	var scopeOptions []huh.Option[string]
	scopeOptions = append(scopeOptions, huh.NewOption("(none)", ""))
	for _, s := range scopes {
		scopeOptions = append(scopeOptions, huh.NewOption(s, s))
	}

	fields := []*huh.Group{}

	fields = append(fields, huh.NewGroup(
		huh.NewConfirm().
			Title("Stage all changes").
			Value(&in.All),
	))

	if in.Conventional {
		fields = append(fields, huh.NewGroup(
			huh.NewSelect[string]().
				Title("Type").
				Options(typeOptions...).
				Value(&in.Type),

			huh.NewSelect[string]().
				Title("Scope").
				Options(scopeOptions...).
				Value(&in.Scope),

			huh.NewInput().
				Title("Summary").
				Value(&in.Summary).
				Placeholder("add user auth"),

			huh.NewText().
				Title("Body").
				Value(&in.Body).
				Placeholder("what changed and why"),

			huh.NewConfirm().
				Title("Breaking change").
				Value(&in.Breaking),
		))
	} else {
		fields = append(fields, huh.NewGroup(
			huh.NewInput().
				Title("Message").
				Value(&in.Summary).
				Placeholder("update stuff"),

			huh.NewText().
				Title("Body").
				Value(&in.Body).
				Placeholder("optional"),
		))
	}

	form := huh.NewForm(fields...)
	if err := form.Run(); err != nil {
		return CommitPromptInput{}, err
	}

	in.Type = strings.TrimSpace(in.Type)
	in.Scope = strings.TrimSpace(in.Scope)
	in.Summary = strings.TrimSpace(in.Summary)
	in.Body = strings.TrimSpace(in.Body)

	if in.Conventional {
		if requireScope && in.Scope == "" {
			return CommitPromptInput{}, huh.ErrUserAborted
		}
	}

	return in, nil
}
