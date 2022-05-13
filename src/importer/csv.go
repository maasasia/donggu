package importer

import (
	"encoding/csv"
	"io"
	"os"
	"path"
	"strings"

	"github.com/maasasia/donggu/dictionary"
	"github.com/pkg/errors"
)

const keyColumnName = "key"
const indexColumnName = "index"
const contextColumnName = "context"

type CsvDictionaryImporter struct{}

func (c CsvDictionaryImporter) OpenMetadataFile(projectRoot string) (io.ReadCloser, error) {
	return nil, nil
}

func (c CsvDictionaryImporter) ImportMetadata(file io.Reader) (dictionary.Metadata, error) {
	return dictionary.Metadata{}, errors.New("unsupported")
}

func (c CsvDictionaryImporter) OpenContentFile(projectRoot string) (io.ReadCloser, error) {
	return os.OpenFile(path.Join(projectRoot, "content.csv"), os.O_RDONLY, 0)
}

func (c CsvDictionaryImporter) ImportContent(file io.Reader, metadata dictionary.Metadata) (dictionary.ContentRepresentation, error) {
	langSet := metadata.SupportedLanguageSet()

	reader := csv.NewReader(file)
	header, headerReadErr := reader.Read()
	if headerReadErr != nil {
		return &dictionary.FlattenedContent{}, errors.Wrap(headerReadErr, "error while reading csv header")
	}
	for index, col := range header {
		col = strings.ToLower(col)
		header[index] = col
		_, isLanguage := langSet[col]
		if !isLanguage && col != contextColumnName && col != indexColumnName && col != keyColumnName {
			return &dictionary.FlattenedContent{}, errors.Errorf("invalid header '%s' at index %d", col, index)
		}
	}

	result := dictionary.FlattenedContent{}
	for {
		line, lineErr := reader.Read()
		if lineErr == io.EOF {
			break
		}
		if lineErr != nil {
			return &dictionary.FlattenedContent{}, errors.Wrap(lineErr, "error while reading csv body")
		}

		keyName := ""
		entry := dictionary.Entry{}

		for index, col := range line {
			colName := header[index]
			if colName == keyColumnName {
				keyName = col
			} else if colName != indexColumnName {
				entry[colName] = col
			}
		}
		if _, ok := result[dictionary.EntryKey(keyName)]; ok {
			return &dictionary.FlattenedContent{}, errors.Errorf("duplicate key '%s'", keyName)
		}
		result[dictionary.EntryKey(keyName)] = entry
	}

	return &result, nil
}
