package importer

import (
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/maasasia/donggu/dictionary"
	"github.com/pkg/errors"
)

type jsonContentType map[string]map[string]string
type jsonMetadataType struct {
	Version            string                            `json:"version"`
	RequiredLanguages  []string                          `json:"required_languages"`
	SupportedLanguages []string                          `json:"supported_languages"`
	ExporterOptions    map[string]map[string]interface{} `json:"exporter_options"`
}

type JsonDictionaryImporter struct{}

func (j JsonDictionaryImporter) OpenMetadataFile(projectRoot string) (io.ReadCloser, error) {
	return os.OpenFile(path.Join(projectRoot, "metadata.json"), os.O_RDONLY, 0)
}

func (j JsonDictionaryImporter) OpenContentFile(projectRoot string) (io.ReadCloser, error) {
	return os.OpenFile(path.Join(projectRoot, "content.json"), os.O_RDONLY, 0)
}

func (j JsonDictionaryImporter) ImportContent(file io.Reader, _ dictionary.Metadata) (dictionary.ContentRepresentation, error) {
	decoder := json.NewDecoder(file)
	decoded := jsonContentType{}
	if err := decoder.Decode(&decoded); err != nil {
		return &dictionary.FlattenedContent{}, errors.Wrap(err, "failed to decode JSON")
	}
	result := dictionary.FlattenedContent{}
	for entryKey, entry := range decoded {
		result[dictionary.EntryKey(entryKey)] = dictionary.Entry(entry)
	}
	return &result, nil
}

func (j JsonDictionaryImporter) ImportMetadata(file io.Reader) (dictionary.Metadata, error) {
	decoder := json.NewDecoder(file)
	decoded := jsonMetadataType{}
	if err := decoder.Decode(&decoded); err != nil {
		return dictionary.Metadata{}, errors.Wrap(err, "failed to decode JSON")
	}

	result := dictionary.Metadata{}
	result.RequiredLanguages = decoded.RequiredLanguages
	result.SupportedLanguages = decoded.SupportedLanguages
	result.Version = decoded.Version
	result.ExporterOptions = decoded.ExporterOptions
	if result.Version == "" {
		return dictionary.Metadata{}, errors.New("version missing")
	}

	return result, nil
}
