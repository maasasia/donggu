package exporter

import (
	"os"
	"path"
	"strings"

	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
	"github.com/maasasia/donggu/exporter/golang"
	"github.com/maasasia/donggu/util"
	"github.com/pkg/errors"
)

// GolangDictionaryExporter is a DictionaryProjectExporter
// generating a Go module for using the dictionary.
type GolangDictionaryExporter struct{}

func (g GolangDictionaryExporter) Export(
	projectRoot string,
	content dictionary.ContentRepresentation,
	metadata dictionary.Metadata,
	options OptionMap,
) error {
	if err := g.prepareProject(projectRoot, options["packageName"].(string)); err != nil {
		return errors.Wrap(err, "failed to prepare project")
	}

	builder := golang.NewGolangBuilder(metadata)
	if err := builder.Run(content.ToTree()); err != nil {
		return err
	}

	return builder.Build(metadata, projectRoot)
}

func (g GolangDictionaryExporter) prepareProject(projectRoot, packageName string) error {
	if err := os.RemoveAll(projectRoot); err != nil {
		return err
	}
	if err := code.CopyTemplateTo("golang", projectRoot, code.CopyTemplateOptions{}); err != nil {
		return errors.Wrap(err, "failed to prepare project")
	}

	packageRenameErr := util.BatchReplaceFiles(
		[]string{path.Join(projectRoot, "go.mod"), path.Join(projectRoot, "donggu.go")},
		"github.com/ghost/donggu",
		packageName,
	)
	if packageRenameErr != nil {
		return errors.Wrap(packageRenameErr, "failed to write package names")
	}

	return nil
}

func (g GolangDictionaryExporter) ValidateOptions(options OptionMap) error {
	convOpts := map[string]interface{}(options)
	if packageName, err := util.SafeAccessMap[string](&convOpts, "packageName"); err == nil {
		if strings.TrimSpace(packageName) == "" {
			return errors.New("package name (key 'packageName') should not be empty")
		}
	} else {
		return err
	}
	return nil
}
