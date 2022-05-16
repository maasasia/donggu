package exporter

import (
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/maasasia/donggu/dictionary"
	"github.com/pkg/errors"
)

// JsonDictionaryExporter is a DictionaryFileExporter.
type JsonDictionaryExporter struct{}

func (j JsonDictionaryExporter) Export(
	projectRoot string,
	content dictionary.ContentRepresentation,
	metadata dictionary.Metadata,
) error {
	metadataFilePath := path.Join(projectRoot, "metadata.json")
	contentFilePath := path.Join(projectRoot, "content.json")

	metadataFile, err := os.OpenFile(metadataFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "failed to open metadata file")
	}
	defer metadataFile.Close()

	contentFile, err := os.OpenFile(contentFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "failed to open content file")
	}
	defer contentFile.Close()

	if err := j.ExportMetadata(metadataFile, metadata); err != nil {
		return errors.Wrapf(err, "failed to write metadata")
	}
	if err := j.ExportContent(contentFile, content, metadata); err != nil {
		return errors.Wrapf(err, "failed to write content")
	}
	return nil
}

func (j JsonDictionaryExporter) ExportContent(
	file io.Writer,
	content dictionary.ContentRepresentation,
	_ dictionary.Metadata,
) error {
	flattened := content.ToFlattened()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(flattened); err != nil {
		return errors.Wrap(err, "failed to encode JSON")
	}
	return nil
}

func (j JsonDictionaryExporter) ExportMetadata(file io.Writer, metadata dictionary.Metadata) error {
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	jsonObj := map[string]interface{}{
		"version":             metadata.Version,
		"required_languages":  metadata.RequiredLanguages,
		"supported_languages": metadata.SupportedLanguages,
	}

	if err := encoder.Encode(jsonObj); err != nil {
		return errors.Wrap(err, "failed to encode JSON")
	}
	return nil
}
