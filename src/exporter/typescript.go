package exporter

import (
	"io/fs"
	"os"
	"path"
	"regexp"

	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
	"github.com/pkg/errors"
)

// TypescriptDictionaryExporter is a DictionaryProjectExporter
// generating a Typescript package for using the dictionary.
type TypescriptDictionaryExporter struct{}

func (t TypescriptDictionaryExporter) Export(
	projectRoot string,
	content dictionary.ContentRepresentation,
	metadata dictionary.Metadata,
	_ OptionMap,
) error {
	if err := t.prepareProject(projectRoot); err != nil {
		return errors.Wrap(err, "failed to prepare project")
	}

	builder := newTypescriptBuilder(metadata)
	if err := builder.Run(content.ToTree()); err != nil {
		return err
	}

	file, err := os.OpenFile(path.Join(projectRoot, "generated/dictionary.ts"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer file.Close()

	builder.Build(metadata, file)
	return nil
}

func (t TypescriptDictionaryExporter) prepareProject(projectRoot string) error {
	if err := os.RemoveAll(projectRoot); err != nil {
		return err
	}

	skipRegex := regexp.MustCompile("(dist|node_modules|generated)")
	skipFunc := func(src string) (bool, error) {
		return skipRegex.MatchString(src), nil
	}
	if err := code.CopyTemplateTo("typescript", projectRoot, code.CopyTemplateOptions{Skip: skipFunc}); err != nil {
		return err
	}
	return os.Mkdir(path.Join(projectRoot, "generated"), fs.ModePerm)
}

func (t TypescriptDictionaryExporter) ValidateOptions(options OptionMap) error {
	return nil
}
