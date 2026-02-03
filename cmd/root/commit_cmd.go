package root

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"gitflow/internal/ui"
	"gitflow/internal/workflow"
)

func commitCmd() *cobra.Command {
	var all bool
	var interactive bool

	var msg string
	var body string

	var ctype string
	var scope string
	var breaking bool

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Create a commit using conventions and optional prompts",
		RunE: func(cmd *cobra.Command, args []string) error {
			common, err := commonFromCmd(cmd)
			if err != nil {
				return err
			}

			repoPath, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			useInteractive := interactive
			if !useInteractive {
				if strings.TrimSpace(msg) == "" && strings.TrimSpace(ctype) == "" && strings.TrimSpace(scope) == "" && strings.TrimSpace(body) == "" && !breaking && !all {
					useInteractive = true
				}
			}

			if useInteractive {
				def := ui.CommitPromptInput{
					Conventional: common.ConfigResult.Config.Commits.Conventional,
					Type:         ctype,
					Scope:        scope,
					Summary:      msg,
					Body:         body,
					Breaking:     breaking,
					All:          all,
				}

				types := common.ConfigResult.Config.Commits.Types
				if len(types) == 0 {
					types = []string{"feat", "fix", "docs", "refactor", "test", "chore"}
				}

				scopes := common.ConfigResult.Config.Commits.Scopes
				in, err := ui.PromptCommit(types, scopes, common.ConfigResult.Config.Commits.RequireScope, def)
				if err != nil {
					return err
				}

				all = in.All
				breaking = in.Breaking
				ctype = in.Type
				scope = in.Scope
				msg = in.Summary
				body = in.Body
			}

			out, err := workflow.Commit(common.ConfigResult.Config, workflow.CommitOptions{
				RepoPath: repoPath,
				All:      all,
				Message:  msg,
				Body:     body,
				Type:     ctype,
				Scope:    scope,
				Breaking: breaking,
			})
			if err != nil {
				return err
			}

			common.UI.Header("Commit created")
			common.UI.Line("Message")
			common.UI.Line(out.Message)

			return nil
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "Stage all changes before commit")
	cmd.Flags().BoolVar(&interactive, "interactive", false, "Prompt for fields")

	cmd.Flags().StringVar(&msg, "message", "", "Commit summary or message")
	cmd.Flags().StringVar(&body, "body", "", "Commit body")

	cmd.Flags().StringVar(&ctype, "type", "", "Conventional type")
	cmd.Flags().StringVar(&scope, "scope", "", "Conventional scope")
	cmd.Flags().BoolVar(&breaking, "breaking", false, "Mark as breaking change")

	return cmd
}
