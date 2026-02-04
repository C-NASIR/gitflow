// Package root wires the CLI commands for gitflow.
package root

import (
	"fmt"
	"gitflow/cmd/root/branch"
	"gitflow/cmd/root/pr"
	"gitflow/cmd/root/provider"
	"gitflow/cmd/root/release"
	"gitflow/internal/cli"
	"gitflow/internal/version"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitflow",
	Short: "getflow automates common git workflows",
	Long:  "gitflow is a professional CLI to automate common git workflow and optionally integrate git hosting providers",
}

var (
	flagNoColor bool
	flagEmoji   bool
	flagVerbose bool
)

// Execute runs the gitflow CLI.
func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	rootCmd.Version = version.String()
	rootCmd.SetVersionTemplate("gitflow {{.Version}}\n")

	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().BoolVar(&flagEmoji, "emoji", false, "Force emoji output")
	rootCmd.PersistentFlags().BoolVar(&flagVerbose, "verbose", false, "Enable verbose output")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		overrides := cli.UIOverrides{}
		if cmd.Flags().Changed("no-color") {
			color := !flagNoColor
			overrides.Color = &color
		}
		if cmd.Flags().Changed("emoji") {
			emoji := flagEmoji
			overrides.Emoji = &emoji
		}
		if cmd.Flags().Changed("verbose") {
			verbose := flagVerbose
			overrides.Verbose = &verbose
		}
		cli.SetUIOverrides(overrides)
		return nil
	}

	rootCmd.AddCommand(VersionCmd())
	rootCmd.AddCommand(statusCmd())
	rootCmd.AddCommand(doctorCmd())
	rootCmd.AddCommand(configCmd())
	rootCmd.AddCommand(startCmd())
	rootCmd.AddCommand(syncCmd())
	rootCmd.AddCommand(cleanupCmd())
	rootCmd.AddCommand(commitCmd())
	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(provider.Cmd())
	rootCmd.AddCommand(pr.Cmd())
	rootCmd.AddCommand(branch.Cmd())
	rootCmd.AddCommand(release.Cmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

}
