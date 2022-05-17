package exporter

import (
	"io"

	"github.com/maasasia/donggu/dictionary"
)

type OptionMap map[string]interface{}

// DictionaryFileExporter exports dictonary content or metadata as a single file.
type DictionaryFileExporter interface {
	ValidateOptions(options OptionMap) error
	ExportContent(file io.Writer, content dictionary.ContentRepresentation, metadata dictionary.Metadata, options OptionMap) error
	ExportMetadata(file io.Writer, metadata dictionary.Metadata, options OptionMap) error
}

// DictionaryProjectExporter exports dictionary content as a collection of files and folders,
// such as a package used in application code or a HTML rendered dictionary visualization.
type DictionaryProjectExporter interface {
	ValidateOptions(options OptionMap) error
	Export(projectRoot string, content dictionary.ContentRepresentation, metadata dictionary.Metadata, options OptionMap) error
}
