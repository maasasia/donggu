package exporter

import (
	"os"

	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
	"github.com/pkg/errors"
)

// GolangDictionaryExporter is a DictionaryProjectExporter
// generating a Go module for using the dictionary.
type GolangDictionaryExporter struct{}

func (g GolangDictionaryExporter) Export(
	projectRoot string,
	content dictionary.ContentRepresentation,
	metadata dictionary.Metadata,
) error {
	if err := g.prepareProject(projectRoot); err != nil {
		return errors.Wrap(err, "failed to prepare project")
	}

	builder := newGolangBuilder(metadata)
	if err := builder.Run(content.ToTree()); err != nil {
		return err
	}

	return builder.Build(metadata, projectRoot)
}

func (g GolangDictionaryExporter) prepareProject(projectRoot string) error {
	if err := os.RemoveAll(projectRoot); err != nil {
		return err
	}

	if err := code.CopyTemplateTo("golang", projectRoot, code.CopyTemplateOptions{}); err != nil {
		return err
	}
	return nil
}
