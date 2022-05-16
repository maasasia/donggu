package exporter

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/maasasia/donggu/dictionary"
	"github.com/pkg/errors"
)

// CsvDictionaryExporter is a DictionaryFileExporter.
type CsvDictionaryExporter struct{}

func (c CsvDictionaryExporter) ExportContent(
	file io.Writer,
	content dictionary.ContentRepresentation,
	metadata dictionary.Metadata,
) error {
	flattened := content.ToFlattened()

	encoder := csv.NewWriter(file)
	defer encoder.Flush()

	locales := make([]string, 0, len(metadata.SupportedLanguages))
	locales = append(locales, metadata.SupportedLanguages...)

	header := make([]string, 0, len(locales)+2)
	header = append(header, "index", "key")
	header = append(header, locales...)
	if err := encoder.Write(header); err != nil {
		return errors.Wrap(err, "error while writing file")
	}

	index := 0
	for key, value := range *flattened {
		index += 1
		row := make([]string, 0, len(locales)+2)
		row = append(row, fmt.Sprintf("%d", index), string(key))
		for _, locale := range locales {
			if localeValue, ok := value[locale]; ok {
				row = append(row, localeValue)
			} else {
				row = append(row, "")
			}
		}
		if err := encoder.Write(row); err != nil {
			return errors.Wrap(err, "error while writing file")
		}
	}

	return nil
}

func (c CsvDictionaryExporter) ExportMetadata(file io.Writer, metadata dictionary.Metadata) error {
	return errors.New("unsupported")
}
