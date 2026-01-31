package root

import (
	"fmt"
	"gitflow/internal/version"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitflow",
	Short: "getflow automates common git workflows",
	Long:  "gitflow is a professional CLI to automate common git workflow and optionally integrate git hosting providers",
}

func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	rootCmd.Version = version.String()
	rootCmd.SetVersionTemplate("gitflow {{.Version}}\n")

	rootCmd.AddCommand(VersionCmd())
	rootCmd.AddCommand(statusCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

}
