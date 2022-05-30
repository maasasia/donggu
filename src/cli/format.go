package cli

import (
	"github.com/maasasia/donggu/dictionary"
	"github.com/maasasia/donggu/exporter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func execFormatCommand(cmd *cobra.Command, _ []string) error {
	projectRoot, err := getProjectRoot(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to resolve project root")
	}
	content, meta, err := loadProject(projectRoot)
	if err != nil {
		return errors.Wrap(err, "failed to load project")
	}

	validateErr := content.Validate(meta, dictionary.ContentValidationOptions{})
	if validateErr != nil {
		return errors.Wrap(validateErr, "content file has errors")
	}

	exportErr := exporter.JsonDictionaryExporter{}.Export(projectRoot, content, meta, exporter.OptionMap{})
	if exportErr != nil {
		return errors.Wrap(err, "failed to save file")
	}
	return nil
}

func initFormatCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "fmt format file",
		Aliases: []string{"format"},
		Short:   "Format content and metadata file",
		Run:     wrapExecCommand(execFormatCommand),
	}
	return cmd
}
