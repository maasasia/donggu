package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func initGenerateCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "generate [format] [path]",
		Short: "Generate something",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			exporterName, targetRoot := args[0], args[1]
			content, meta, err := loadProjectFromCommand(cmd)
			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}

			targetRoot, err = filepath.Abs(targetRoot)
			if err != nil {
				fmt.Printf("Invalid target path '%s'.\n%s\n", args[1], err)
				os.Exit(1)
			}

			if exporter := loadProjectExporter(exporterName); exporter != nil {
				err := exporter.Export(targetRoot, content, meta)
				if err != nil {
					fmt.Printf("Failed to export project.\n%s\n", err)
					os.Exit(1)
				}
			} else if exporter := loadFileExporter(exporterName); exporter != nil {
				file, err := os.OpenFile(targetRoot, os.O_TRUNC, os.ModePerm)
				if err != nil {
					fmt.Printf("Failed to open target file.\n%s\n", err)
					os.Exit(1)
				}
				defer file.Close()

				err = exporter.ExportContent(file, content, meta)
				if err != nil {
					fmt.Printf("Failed to export project.\n%s\n", err)
					os.Exit(1)
				}
			} else {
				fmt.Printf("Unknown exporter '%s'\n", exporter)
				os.Exit(1)
			}
		},
	}

	return cmd
}
