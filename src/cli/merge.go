package cli

import (
	"path/filepath"

	"github.com/maasasia/donggu/dictionary"
	"github.com/maasasia/donggu/exporter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func execMergeCommand(cmd *cobra.Command, args []string) error {
	format, filePath := args[0], args[1]
	projectRoot, err := getProjectRoot(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to resolve project root")
	}
	content, meta, err := loadProject(projectRoot)
	if err != nil {
		return errors.Wrap(err, "failed to load project")
	}

	filePath, err = filepath.Abs(filePath)
	if err != nil {
		return errors.Wrapf(err, "invalid target path '%s'", filePath)
	}

	mergeContent, err := readContentFile(meta, format, filePath)
	if err != nil {
		return errors.Wrap(err, "failed to merging file")
	}

	content = dictionary.MergeContent(mergeContent, content)
	exportErr := exporter.JsonDictionaryExporter{}.Export(projectRoot, content, meta, exporter.OptionMap{})
	if exportErr != nil {
		return errors.Wrap(err, "failed to save merged file")
	}
	return nil
}

func initMergeCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "merge format file",
		Short: "Merge a content file to the current project",
		Args:  cobra.ExactArgs(2),
		Run:   wrapExecCommand(execMergeCommand),
	}
	return cmd
}
