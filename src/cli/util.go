package cli

import (
	"os"
	"path/filepath"

	"github.com/maasasia/donggu/dictionary"
	"github.com/maasasia/donggu/importer"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func loadProjectFromCommand(cmd *cobra.Command) (content dictionary.ContentRepresentation, meta dictionary.Metadata, err error) {
	projectRoot, err := getProjectRoot(cmd)
	if err != nil {
		err = errors.Wrap(err, "failed to resolve project root")
		return
	}
	content, meta, err = loadProject(projectRoot)
	if err != nil {
		err = errors.Wrap(err, "Failed to load project")
	}
	return
}

func loadProject(projectRoot string) (content dictionary.ContentRepresentation, meta dictionary.Metadata, err error) {
	jsonImporter := importer.JsonDictionaryImporter{}

	metaFile, err := jsonImporter.OpenMetadataFile(projectRoot)
	if err != nil {
		err = errors.Wrap(err, "failed to open metadata file")
		return
	}
	defer metaFile.Close()

	contentFile, err := jsonImporter.OpenContentFile(projectRoot)
	if err != nil {
		err = errors.Wrap(err, "failed to open content file")
		return
	}
	defer contentFile.Close()

	meta, err = jsonImporter.ImportMetadata(metaFile)
	if err != nil {
		err = errors.Wrap(err, "failed to read metadata file")
		return
	}
	content, err = jsonImporter.ImportContent(contentFile, meta)
	if err != nil {
		err = errors.Wrap(err, "failed to read content file")
		return
	}
	return
}

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
