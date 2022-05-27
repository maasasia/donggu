package exporter

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
	"github.com/maasasia/donggu/exporter/typescript"
	"github.com/maasasia/donggu/util"
	"github.com/pkg/errors"
)

// TypescriptReactDictionaryExporter is a DictionaryProjectExporter
// generating a Typescript React package for using the dictionary.
type TypescriptReactDictionaryExporter struct{}

func (t TypescriptReactDictionaryExporter) Export(
	projectRoot string,
	content dictionary.ContentRepresentation,
	metadata dictionary.Metadata,
	options OptionMap,
) error {
	if err := t.prepareProject(projectRoot, metadata, options); err != nil {
		return errors.Wrap(err, "failed to prepare project")
	}

	builder := typescript.NewTypescriptBuilder(metadata, &typescript.ReactBuilderOptions{})
	if err := builder.Run(content.ToTree()); err != nil {
		return err
	}

	file, err := os.OpenFile(path.Join(projectRoot, "generated/dictionary.tsx"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer file.Close()

	builder.Build(metadata, file)
	return nil
}

func (t TypescriptReactDictionaryExporter) prepareProject(projectRoot string, metadata dictionary.Metadata, options OptionMap) error {
	if err := os.RemoveAll(projectRoot); err != nil {
		return err
	}

	skipRegex := regexp.MustCompile("(dist|node_module|generated)")
	skipFunc := func(src string) (bool, error) {
		return skipRegex.MatchString(src), nil
	}
	if err := code.CopyTemplateTo("ts-react", projectRoot, code.CopyTemplateOptions{Skip: skipFunc}); err != nil {
		return errors.Wrap(err, "failed to prepare export project")
	}
	if err := os.Mkdir(path.Join(projectRoot, "generated"), fs.ModePerm); err != nil {
		return errors.Wrap(err, "failed to create generated folder")
	}

	replaceErr := util.MultiReplaceFile(path.Join(projectRoot, "package.json"), []util.ReplaceSet{
		{From: regexp.MustCompile("donggu-template-ts"), To: options["packageName"].(string)},
		{From: regexp.MustCompile(`"version": "1.0.0"`), To: fmt.Sprintf(`"version": "%s"`, metadata.Version)},
	})
	if replaceErr != nil {
		return errors.Wrap(replaceErr, "failed to edit package.json")
	}
	return nil
}

func (t TypescriptReactDictionaryExporter) ValidateOptions(options OptionMap) error {
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
