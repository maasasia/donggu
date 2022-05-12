package exporter

import (
	"encoding/json"
	"os"

	"github.com/maasasia/donggu/dictionary"
	"github.com/pkg/errors"
)

type JsonDictionaryExporter struct{}

func (j JsonDictionaryExporter) ExportContent(filePath string, content dictionary.ContentRepresentation) error {
	flattened := content.ToFlattened()

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(flattened); err != nil {
		return errors.Wrap(err, "failed to encode JSON")
	}
	return nil
}

func (j JsonDictionaryExporter) ExportMetadata(filePath string, metadata dictionary.Metadata) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer file.Close()

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
