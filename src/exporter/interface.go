package exporter

import (
	"io"

	"github.com/maasasia/donggu/dictionary"
)

// DictionaryFileExporter exports dictonary content or metadata as a single file.
type DictionaryFileExporter interface {
	ExportContent(file io.Writer, content dictionary.ContentRepresentation, metadata dictionary.Metadata) error
	ExportMetadata(file io.Writer, metadata dictionary.Metadata) error
}

// DictionaryProjectExporter exports dictionary content as a collection of files and folders,
// such as a package used in application code or a HTML rendered dictionary visualization.
type DictionaryProjectExporter interface {
	Export(projectRoot string, content dictionary.ContentRepresentation, metadata dictionary.Metadata) error
}
