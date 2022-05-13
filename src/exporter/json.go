package exporter

import (
	"encoding/json"
	"io"

	"github.com/maasasia/donggu/dictionary"
	"github.com/pkg/errors"
)

// JsonDictionaryExporter is a DictionaryFileExporter.
type JsonDictionaryExporter struct{}

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

	jsonObj := map[string][]string{
		"required_languages":  metadata.RequiredLanguages,
		"supported_languages": metadata.SupportedLanguages,
	}

	if err := encoder.Encode(jsonObj); err != nil {
		return errors.Wrap(err, "failed to encode JSON")
	}
	return nil
}
