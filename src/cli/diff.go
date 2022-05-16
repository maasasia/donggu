package cli

import (
	"path/filepath"

	"github.com/maasasia/donggu/dictionary"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const diffCommandDescription = `
diff compares the difference from a content file to the current project's content.json file.

The source is the content file and the destination is the content.json file.
For instance, if a key 'screens.login' exists in the content file
but is deleted in content.json, 'screens.login' would be marked as deleted.
This behavior can be reversed by the '--reverse' flag.`

func execDiffCommand(cmd *cobra.Command, args []string) error {
	reverse, _ := cmd.Flags().GetBool("reverse")
	otherFileFormat, filePath := args[0], args[1]
	content, meta, err := loadProjectFromCommand(cmd)
	if err != nil {
		return err
	}

	filePath, err = filepath.Abs(filePath)
	if err != nil {
		return errors.Wrapf(err, "invalid target path '%s'", filePath)
	}

	otherFile, err := readContentFile(meta, otherFileFormat, filePath)
	if err != nil {
		return errors.Wrap(err, "failed to read other file")
	}

	var diff dictionary.ContentDifference
	if reverse {
		diff = dictionary.DiffContents(otherFile, content)
	} else {
		diff = dictionary.DiffContents(content, otherFile)
	}

	var writeErr error
	if useCsv, _ := cmd.Flags().GetBool("csv"); useCsv {
		writeErr = outputDiffToCsv(diff)
	} else {
		outputDiffToConsole(diff, reverse)
		writeErr = nil
	}
	if writeErr != nil {
		return errors.Wrap(err, "failed to write output")
	}
	return nil
}

func initDiffCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "diff [--console | --csv] [--reverse] format file",
		Short: "Show differences of a content file against the current project",
		Long:  diffCommandDescription,
		Args:  cobra.ExactArgs(2),
		Run:   wrapExecCommand(execDiffCommand),
	}

	cmd.PersistentFlags().Bool("console", false, "Print the difference to the console")
	cmd.PersistentFlags().Bool("csv", false, "Print the difference to a CSV file")
	cmd.PersistentFlags().Bool("reverse", false, "Reverse the source/destination between files. See --help for details.")

	return cmd
}
