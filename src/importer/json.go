package importer

import (
	"encoding/json"
	"os"
	"path"

	"github.com/maasasia/donggu/dictionary"
	"github.com/pkg/errors"
)

type jsonContentType map[string]map[string]string
type jsonMetadataType map[string][]string

type JsonDictionaryImporter struct{}

func (j JsonDictionaryImporter) ImportContent(filePath string, _ dictionary.Metadata) (dictionary.ContentRepresentation, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		return &dictionary.FlattenedContent{}, errors.Wrap(err, "failed to open file")
	}
	defer file.Close()

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

func (j JsonDictionaryImporter) ImportMetadata(filePath string) (dictionary.Metadata, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		return dictionary.Metadata{}, errors.Wrap(err, "failed to open file")
	}
	defer file.Close()

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
		return dictionary.Metadata{}, errors.Wrap(err, "RequiredLanguages is required but not given")
	}
	supportedLangs, ok := decoded["supported_languages"]
	if ok {
		result.SupportedLanguages = supportedLangs
	} else {
		return dictionary.Metadata{}, errors.Wrap(err, "SupportedLanguages is required but not given")
	}

	return result, nil
}

func (j JsonDictionaryImporter) ResolveProject(projectPath string) (ResolveProjectResult, error) {
	return ResolveProjectResult{
		ContentPath:  path.Join(projectPath, "content.json"),
		MetadataPath: path.Join(projectPath, "metadata.json"),
	}, nil
}
