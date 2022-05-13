package cli

import (
	"github.com/maasasia/donggu/dictionary"
	"github.com/maasasia/donggu/importer"
	"github.com/pkg/errors"
)

func loadProject(projectRoot string) (content dictionary.ContentRepresentation, meta dictionary.Metadata, err error) {
	jsonImporter := importer.JsonDictionaryImporter{}

	metaFile, err := jsonImporter.OpenMetadataFile(projectRoot)
	if err != nil {
		err = errors.Wrap(err, "failed to open metadata file")
		return
	}
	defer metaFile.Close()

	contentFile, err := jsonImporter.OpenContentFile(projectRoot)
	if err != nil {
		err = errors.Wrap(err, "failed to open content file")
		return
	}
	defer contentFile.Close()

	meta, err = jsonImporter.ImportMetadata(metaFile)
	if err != nil {
		err = errors.Wrap(err, "failed to read metadata file")
		return
	}
	content, err = jsonImporter.ImportContent(contentFile, meta)
	if err != nil {
		err = errors.Wrap(err, "failed to read content file")
		return
	}
	return
}
