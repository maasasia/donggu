package cli

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func getProjectRoot(cmd *cobra.Command) (string, error) {
	projectRoot, _ := cmd.Flags().GetString("project")
	if projectRoot == "" {
		workingDir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		projectRoot = workingDir
	} else {
		absPath, err := filepath.Abs(projectRoot)
		if err != nil {
			return "", err
		}
		projectRoot = absPath
	}
	return projectRoot, nil
}
