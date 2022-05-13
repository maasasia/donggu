package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/maasasia/donggu/exporter"
	"github.com/spf13/cobra"
)

var fileExporters = map[string]exporter.DictionaryFileExporter{
	"json": exporter.JsonDictionaryExporter{},
}

var projectExporters = map[string]exporter.DictionaryProjectExporter{
	"typescript": exporter.TypescriptDictionaryExporter{},
}

func addGenerateCommand(parent *cobra.Command) {
	var cmd = &cobra.Command{
		Use:   "generate [format] [path]",
		Short: "Generate something",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			projectRoot, err := getProjectRoot(cmd)
			if err != nil {
				fmt.Printf("Failed to load project root.\n%s\n", err)
				os.Exit(1)
			}
			targetRoot, err := filepath.Abs(args[1])
			if err != nil {
				fmt.Printf("Invalid target path '%s'.\n%s\n", args[1], err)
				os.Exit(1)
			}

			exporter := args[0]
			projectExporter, isProjectExporter := projectExporters[exporter]
			fileExporter, isFileExporter := fileExporters[exporter]
			if !isProjectExporter && !isFileExporter {
				fmt.Printf("Unknown exporter '%s'\n", exporter)
				os.Exit(1)
			}

			content, meta, err := loadProject(projectRoot)
			if err != nil {
				fmt.Printf("Failed to load project.\n%s\n", err)
				os.Exit(1)
			}

			if isProjectExporter {
				err := projectExporter.Export(targetRoot, content, meta)
				if err != nil {
					fmt.Printf("Failed to export project.\n%s\n", err)
					os.Exit(1)
				}
				os.Exit(0)
			}

			file, err := os.OpenFile(targetRoot, os.O_TRUNC, os.ModePerm)
			if err != nil {
				fmt.Printf("Failed to open target file.\n%s\n", err)
				os.Exit(1)
			}
			defer file.Close()

			err = fileExporter.ExportContent(file, content, meta)
			if err != nil {
				fmt.Printf("Failed to export project.\n%s\n", err)
				os.Exit(1)
			}
		},
	}
	parent.AddCommand(cmd)
}
