package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/maasasia/donggu/dictionary"
	"github.com/maasasia/donggu/importer"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func wrapExecCommand(exec func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		startTime := time.Now()
		if err := exec(cmd, args); err != nil {
			fmt.Printf("🚫 Failed to run command.\n%s\n", err)
			os.Exit(1)
		} else {
			duration := time.Since(startTime).Seconds()
			fmt.Printf("✅ Done in %.3fs\n", duration)
		}
	}
}

func loadProjectFromCommand(cmd *cobra.Command) (content dictionary.ContentRepresentation, meta dictionary.Metadata, err error) {
	projectRoot, err := getProjectRoot(cmd)
	if err != nil {
		err = errors.Wrap(err, "failed to resolve project root")
		return
	}
	content, meta, err = loadProject(projectRoot)
	if err != nil {
		err = errors.Wrap(err, "failed to load project")
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

func readContentFile(metadata dictionary.Metadata, format, filePath string) (dictionary.ContentRepresentation, error) {
	isDirectory, err := isPathDirectory(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read content file or directory")
	}

	var importer importer.DictionaryFileImporter
	var file io.ReadCloser

	if isDirectory {
		fullImporter := loadImporter(format)
		if importer == nil {
			return nil, errors.Errorf("unknown import format '%s'", format)
		}
		file, err = fullImporter.OpenContentFile(filePath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open content file")
		}
		importer = fullImporter

	} else {
		importer = loadFileImporter(format)
		if importer == nil {
			return nil, errors.Errorf("unknown import file format '%s'", format)
		}
		file, err = os.OpenFile(filePath, os.O_RDONLY, 0)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open file")
		}
	}
	defer file.Close()
	content, err := importer.ImportContent(file, metadata)
	if err != nil {
		return nil, errors.Wrap(err, "error while parsing file")
	}
	return content, nil
}

func isPathDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, errors.Wrap(err, "cannot stat path")
	}
	return fileInfo.IsDir(), nil
}
