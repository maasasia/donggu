package cli

import (
	"os"
	"path/filepath"

	"github.com/maasasia/donggu/dictionary"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func execExportCommand(cmd *cobra.Command, args []string) error {
	exporterName, targetRoot := args[0], args[1]
	content, meta, err := loadProjectFromCommand(cmd)
	if err != nil {
		return err
	}

	validateErr := meta.Validate()
	if validateErr != nil {
		return errors.Wrap(validateErr, "metadata file has errors")
	}
	validateErr = content.Validate(meta, dictionary.ContentValidationOptions{})
	if validateErr != nil {
		return errors.Wrap(validateErr, "content file has errors")
	}

	targetRoot, err = filepath.Abs(targetRoot)
	if err != nil {
		return errors.Wrapf(err, "invalid target path '%s'", targetRoot)
	}

	exporterOptions := meta.ExporterOption(exporterName)
	if exporter := loadProjectExporter(exporterName); exporter != nil {
		if err := exporter.ValidateOptions(exporterOptions); err != nil {
			return errors.Wrap(err, "invalid options")
		}
		err := exporter.Export(targetRoot, content, meta, exporterOptions)
		if err != nil {
			return errors.Wrap(err, "failed to export project")
		}
	} else if exporter := loadFileExporter(exporterName); exporter != nil {
		if err := exporter.ValidateOptions(exporterOptions); err != nil {
			return errors.Wrap(err, "invalid options")
		}

		file, err := os.OpenFile(targetRoot, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "failed to open target file")
		}
		defer file.Close()

		err = exporter.ExportContent(file, content, meta, exporterOptions)
		if err != nil {
			return errors.Wrap(err, "failed to export project")
		}
	} else {
		return errors.Errorf("Unknown exporter '%s'", exporterName)
	}
	return nil
}

func initExportCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "export format path",
		Short: "Export something",
		Args:  cobra.ExactArgs(2),
		Run:   wrapExecCommand(execExportCommand),
	}

	return cmd
}
