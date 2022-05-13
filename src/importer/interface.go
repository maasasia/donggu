package importer

import (
	"io"

	"github.com/maasasia/donggu/dictionary"
)

type ResolveProjectResult struct {
	ContentPath  string
	MetadataPath string
}

type DictionaryImporter interface {
	DictionaryFileImporter
	DictionaryProjectImporter
}

type DictionaryFileImporter interface {
	ImportContent(r io.Reader, metadata dictionary.Metadata) (dictionary.ContentRepresentation, error)
	ImportMetadata(r io.Reader) (dictionary.Metadata, error)
}

type DictionaryProjectImporter interface {
	// OpenMetadataFile creates a ReadCloser that reads metadata from the given project.
	// Usually this would be an os.File.
	//
	// Returns (nil, nil) if the importer does not support reading metadata.
	OpenMetadataFile(projectRoot string) (io.ReadCloser, error)

	// OpenContentFile creates a ReadCloser that reads content from the given project.
	// Usually this would be an os.File.
	//
	// Returns (nil, nil) if the importer does not support reading content.
	OpenContentFile(projectRoot string) (io.ReadCloser, error)
}
