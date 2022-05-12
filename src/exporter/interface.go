package exporter

import "github.com/maasasia/donggu/dictionary"

type DictionaryExporter interface {
	ExportContent(filePath string, content dictionary.ContentRepresentation) error
	ExportMetadata(filePath string, metadata dictionary.Metadata) error
}
