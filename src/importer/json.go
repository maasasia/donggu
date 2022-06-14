package importer

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/maasasia/donggu/dictionary"
	"github.com/pkg/errors"
)

// Regular expression for validating plural operators
var pluralOperatorRegex = regexp.MustCompile(`^([<>]=?)|(?:(%|\/)([1-9]\d*))$`)

type jsonPluralDefinition struct {
	Op    string `json:"op"`
	Value int    `json:"value"`
}

func (j jsonPluralDefinition) Parse() (dictionary.PluralDefinition, error) {
	match := pluralOperatorRegex.FindAllStringSubmatch(j.Op, -1)
	if len(match) == 0 {
		return dictionary.PluralDefinition{}, errors.Errorf("'%s' is an invalid operator", j.Op)
	}
	if match[0][1] == "" {
		operandStr := match[0][3]
		operand, err := strconv.ParseInt(operandStr, 10, 32)
		if err != nil {
			return dictionary.PluralDefinition{}, errors.Errorf("'%s' is an invalid operand", operandStr)
		}
		return dictionary.PluralDefinition{
			Op:      match[0][2],
			Operand: int(operand),
			Equals:  j.Value,
		}, nil
	} else {
		return dictionary.PluralDefinition{
			Op:     match[0][1],
			Equals: j.Value,
		}, nil
	}
}

type jsonContentType map[string]map[string]string
type jsonMetadataType struct {
	Version            string                            `json:"version"`
	RequiredLanguages  []string                          `json:"required_languages"`
	SupportedLanguages []string                          `json:"supported_languages"`
	ExporterOptions    map[string]map[string]interface{} `json:"exporter_options"`
	Plurals            map[string][]jsonPluralDefinition `json:"plurals"`
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
	if result.Version == "" {
		return dictionary.Metadata{}, errors.New("version missing")
	}
	result.ExporterOptions = decoded.ExporterOptions
	result.Plurals = map[string][]dictionary.PluralDefinition{}
	for lang, defs := range decoded.Plurals {
		result.Plurals[lang] = make([]dictionary.PluralDefinition, 0, len(defs))
		for index, def := range defs {
			converted, err := def.Parse()
			if err != nil {
				return dictionary.Metadata{}, errors.Wrapf(err, "invalid plural definition for '%s' [%d]", lang, index)
			}
			result.Plurals[lang] = append(result.Plurals[lang], converted)
		}
	}
	return result, nil
}
