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
type jsonMetadataType map[string][]string

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
	requiredLangs, ok := decoded["required_languages"]
	if ok {
		result.RequiredLanguages = requiredLangs
	} else {
		return dictionary.Metadata{}, errors.New("RequiredLanguages is required but not given")
	}
	supportedLangs, ok := decoded["supported_languages"]
	if ok {
		result.SupportedLanguages = supportedLangs
	} else {
		return dictionary.Metadata{}, errors.New("SupportedLanguages is required but not given")
	}

	return result, nil
}
