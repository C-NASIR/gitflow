package pr

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"gitflow/internal/config"
	"gitflow/internal/ui"
	"gitflow/internal/workflow"
)

func createCmd() *cobra.Command {
	var title string
	var body string
	var base string
	var draft bool
	var draftSet bool
	var reviewers string
	var labels string
	var remote string
	var interactive bool
	var open bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a pull request for the current branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoPath, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			res, err := config.Load()
			if err != nil {
				return err
			}

			useInteractive := interactive
			if !useInteractive {
				if strings.TrimSpace(title) == "" && strings.TrimSpace(body) == "" && !draftSet && strings.TrimSpace(reviewers) == "" && strings.TrimSpace(labels) == "" && strings.TrimSpace(base) == "" {
					useInteractive = true
				}
			}

			var draftPtr *bool
			if draftSet {
				v := draft
				draftPtr = &v
			}

			if useInteractive {
				def := ui.PRPromptInput{
					Title:       title,
					Description: body,
					Draft:       res.Config.Workflows.PR.Draft,
					Reviewers:   strings.Join(res.Config.Workflows.PR.DefaultReviewers, ","),
					Labels:      strings.Join(res.Config.Workflows.PR.Labels, ","),
					BaseBranch:  base,
					OpenBrowser: open,
				}
				if draftSet {
					def.Draft = draft
				}
				if strings.TrimSpace(reviewers) != "" {
					def.Reviewers = reviewers
				}
				if strings.TrimSpace(labels) != "" {
					def.Labels = labels
				}

				in, err := ui.PromptPR(def)
				if err != nil {
					return err
				}

				title = in.Title
				body = in.Description
				base = in.BaseBranch
				open = in.OpenBrowser

				v := in.Draft
				draftPtr = &v

				reviewers = in.Reviewers
				labels = in.Labels
			}

			out, err := workflow.CreatePR(res.Config, workflow.PRCreateOptions{
				RepoPath:    repoPath,
				Remote:      remote,
				Title:       title,
				Description: body,
				BaseBranch:  base,
				Draft:       draftPtr,
				Reviewers:   splitCSV(reviewers),
				Labels:      splitCSV(labels),
			})
			if err != nil {
				return err
			}

			pr := out.PR
			cmd.Printf("PR created: #%d\n", pr.Number)
			cmd.Printf("Title: %s\n", pr.Title)
			cmd.Printf("URL: %s\n", pr.URL)

			if open {
				_ = ui.OpenURL(pr.URL)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "PR title")
	cmd.Flags().StringVar(&body, "body", "", "PR description")
	cmd.Flags().StringVar(&base, "base", "", "Base branch")
	cmd.Flags().StringVar(&remote, "remote", "origin", "Remote name")
	cmd.Flags().StringVar(&reviewers, "reviewers", "", "Comma separated reviewers")
	cmd.Flags().StringVar(&labels, "labels", "", "Comma separated labels")
	cmd.Flags().BoolVar(&interactive, "interactive", false, "Prompt for missing fields")
	cmd.Flags().BoolVar(&open, "open", false, "Open PR in browser after creation")
	cmd.Flags().BoolVar(&draft, "draft", false, "Create as draft PR")

	cmd.Flags().Lookup("draft").NoOptDefVal = "true"
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("draft") {
			draftSet = true
		}
	}

	return cmd
}

func splitCSV(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}
