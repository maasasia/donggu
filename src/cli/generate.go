package cli

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func execGenerateCommand(cmd *cobra.Command, args []string) error {
	exporterName, targetRoot := args[0], args[1]
	content, meta, err := loadProjectFromCommand(cmd)
	if err != nil {
		return err
	}

	targetRoot, err = filepath.Abs(targetRoot)
	if err != nil {
		return errors.Wrapf(err, "invalid target path '%s'", targetRoot)
	}

	if exporter := loadProjectExporter(exporterName); exporter != nil {
		err := exporter.Export(targetRoot, content, meta)
		if err != nil {
			return errors.Wrap(err, "failed to export project")
		}
	} else if exporter := loadFileExporter(exporterName); exporter != nil {
		file, err := os.OpenFile(targetRoot, os.O_TRUNC, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "failed to open target file")
		}
		defer file.Close()

		err = exporter.ExportContent(file, content, meta)
		if err != nil {
			return errors.Wrap(err, "failed to export project")
		}
	} else {
		return errors.Errorf("Unknown exporter '%s'", exporterName)
	}
	return nil
}

func initGenerateCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "generate [format] [path]",
		Short: "Generate something",
		Args:  cobra.ExactArgs(2),
		Run:   wrapExecCommand(execGenerateCommand),
	}

	return cmd
}
