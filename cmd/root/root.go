// Package root wires the CLI commands for gitflow.
package root

import (
	"fmt"
	"gitflow/cmd/root/pr"
	"gitflow/cmd/root/provider"
	"gitflow/internal/version"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitflow",
	Short: "getflow automates common git workflows",
	Long:  "gitflow is a professional CLI to automate common git workflow and optionally integrate git hosting providers",
}

// Execute runs the gitflow CLI.
func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	rootCmd.Version = version.String()
	rootCmd.SetVersionTemplate("gitflow {{.Version}}\n")

	rootCmd.AddCommand(VersionCmd())
	rootCmd.AddCommand(statusCmd())
	rootCmd.AddCommand(configCmd())
	rootCmd.AddCommand(startCmd())
	rootCmd.AddCommand(syncCmd())
	rootCmd.AddCommand(cleanupCmd())
	rootCmd.AddCommand(commitCmd())
	rootCmd.AddCommand(provider.Cmd())
	rootCmd.AddCommand(pr.Cmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

}
