package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "donggu",
	Short: "Donggu is a simple cli for managing i18n text data",
}

func init() {
	rootCmd.PersistentFlags().StringP("project", "P", "", "Project folder (default: current directory)")
	rootCmd.AddCommand(initExportCommand())
	rootCmd.AddCommand(initMergeCommand())
	rootCmd.AddCommand(initFormatCommand())
	rootCmd.AddCommand(initDiffCommand())
	rootCmd.AddCommand(initInitCommand())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
