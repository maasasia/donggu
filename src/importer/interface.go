package importer

import "github.com/maasasia/donggu/dictionary"

type ResolveProjectResult struct {
	ContentPath  string
	MetadataPath string
}

type DictionaryImporter interface {
	ImportContent(filePath string, metadata dictionary.Metadata) (dictionary.ContentRepresentation, error)
	ImportMetadata(filePath string) (dictionary.Metadata, error)
	ResolveProject(projectPath string) (ResolveProjectResult, error)
}
