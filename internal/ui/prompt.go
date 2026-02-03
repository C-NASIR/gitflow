package ui

import (
	"strings"

	"github.com/charmbracelet/huh"
)

type PRPromptInput struct {
	Title       string
	Description string
	Draft       bool
	Reviewers   string
	Labels      string
	BaseBranch  string
	OpenBrowser bool
}

func PromptPR(defaults PRPromptInput) (PRPromptInput, error) {
	in := defaults

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title").
				Value(&in.Title).
				Placeholder("Add auth"),

			huh.NewText().
				Title("Description").
				Value(&in.Description).
				Placeholder("What changed and why"),

			huh.NewConfirm().
				Title("Draft").
				Value(&in.Draft),

			huh.NewInput().
				Title("Reviewers (comma separated)").
				Value(&in.Reviewers).
				Placeholder("alice,bob"),

			huh.NewInput().
				Title("Labels (comma separated)").
				Value(&in.Labels).
				Placeholder("needs-review"),

			huh.NewInput().
				Title("Base branch").
				Value(&in.BaseBranch).
				Placeholder("main"),

			huh.NewConfirm().
				Title("Open in browser after create").
				Value(&in.OpenBrowser),
		),
	)

	err := form.Run()
	if err != nil {
		return PRPromptInput{}, err
	}

	in.Title = strings.TrimSpace(in.Title)
	in.Description = strings.TrimSpace(in.Description)
	in.Reviewers = strings.TrimSpace(in.Reviewers)
	in.Labels = strings.TrimSpace(in.Labels)
	in.BaseBranch = strings.TrimSpace(in.BaseBranch)

	return in, nil
}
