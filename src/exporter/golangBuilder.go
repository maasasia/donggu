package exporter

import (
	"os"
	"path"
	"time"

	"github.com/dave/jennifer/jen"
	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
	"github.com/pkg/errors"
)

type golangBuilder struct {
	builder          *golangCodeBuilder
	contentValidator dictionary.ContentValidator
}

func newGolangBuilder(metadata dictionary.Metadata) *golangBuilder {
	return &golangBuilder{
		builder: newGolangCodeBuilder(),
		contentValidator: dictionary.NewContentValidator(metadata, dictionary.ContentValidationOptions{
			SkipLangSupportCheck: true,
		}),
	}
}

func (g *golangBuilder) Build(metadata dictionary.Metadata, projectRoot string) error {
	now := time.Now()

	operations := map[string]func(f *os.File) error{
		"data.go": func(f *os.File) error {
			return g.builder.outputDataFile(f, now)
		},
		"language.go": func(f *os.File) error {
			return g.builder.outputLanguageFile(f, metadata, now)
		},
		"nodes.go": func(f *os.File) error {
			return g.builder.outputNodeFile(f, now)
		},
	}

	for filename, saveFile := range operations {
		file, err := g.openFile(projectRoot, filename)
		if err != nil {
			return errors.Wrap(err, "build failed")
		}
		err = saveFile(file)
		file.Close()
		if err != nil {
			return errors.Wrap(err, "build failed")
		}
	}
	return nil
}

func (g *golangBuilder) openFile(projectRoot, filename string) (*os.File, error) {
	f, err := os.OpenFile(
		path.Join(projectRoot, "generated", filename),
		os.O_CREATE|os.O_TRUNC|os.O_RDWR,
		os.ModePerm,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open '%s'", filename)
	}
	return f, err
}

func (g *golangBuilder) Run(content *dictionary.ContentNode) error {
	return g.walk(content, dictionary.EntryKey(""), nil, 0)
}

func (g *golangBuilder) walk(
	contentNode *dictionary.ContentNode,
	positionKey dictionary.EntryKey,
	selfNameEntry dictionary.Entry,
	depth int,
) error {
	childPropertyNames := map[string]struct{}{}

	entriesToSkip := map[string]struct{}{}

	for key, child := range contentNode.Children {
		propertyName := code.ToCamelCase(key)
		childPropertyNames[propertyName] = struct{}{}

		var err error
		if _, ok := contentNode.Entries[key]; ok {
			entriesToSkip[key] = struct{}{}
			err = g.walk(child, positionKey.NewChild(key), contentNode.Entries[key], depth+1)
		} else {
			err = g.walk(child, positionKey.NewChild(key), nil, depth+1)
		}
		if err != nil {
			return err
		}
	}
	for key, entry := range contentNode.Entries {
		if _, ok := entriesToSkip[key]; ok {
			continue
		}
		err := g.addEntry(entry, positionKey.NewChild(key), false)
		if err != nil {
			return err
		}
	}
	if selfNameEntry != nil {
		err := g.addEntry(selfNameEntry, positionKey.NewChild("_"), true)
		if err != nil {
			return err
		}
	}

	g.writeNodeToBuilder(positionKey, &childPropertyNames, depth == 0)
	return nil
}

func (g *golangBuilder) addEntry(entry dictionary.Entry, entryKey dictionary.EntryKey, isSelfEntry bool) (err error) {
	templateKeys, validateErr := g.contentValidator.Validate(entry)
	if validateErr != nil {
		err = errors.Wrap(validateErr, "failed to add leaf")
		return
	}
	callArgs, paramArgs := g.buildEntryArgumentBlock(templateKeys)
	g.builder.writeEntryType(entryKey, paramArgs)
	g.builder.writeEntryMethod(entryKey, paramArgs, callArgs)
	for lang := range entry {
		if lang == "context" {
			continue
		}
		g.builder.writeEntryImpl(entryKey, lang, paramArgs, g.buildFormatterReturnValue(entry, lang, templateKeys))
	}
	return
}

func (g *golangBuilder) buildEntryArgumentBlock(argTypes map[string]dictionary.TemplateKeyFormat) (callArgs, paramArgs []jen.Code) {
	callArgs = make([]jen.Code, 0, len(argTypes))
	paramArgs = make([]jen.Code, 0, len(argTypes))
	for k, v := range argTypes {
		callArg, paramArg := g.resolveEntryArgumentParam(code.TemplateKeyToCamelCase(k), v)
		callArgs = append(callArgs, callArg)
		paramArgs = append(paramArgs, paramArg)
	}
	return
}
func (g *golangBuilder) writeNodeToBuilder(
	parentKey dictionary.EntryKey,
	childPropertyNames *map[string]struct{},
	isRoot bool,
) {
	if !isRoot {
		g.builder.writeNodeType(parentKey)
	}
	for childName := range *childPropertyNames {
		g.builder.writeNodeChild(parentKey, childName, isRoot)
	}
}

func (g *golangBuilder) buildFormatterReturnValue(
	entry dictionary.Entry,
	lang string,
	argTypes map[string]dictionary.TemplateKeyFormat,
) *jen.Statement {
	params := make([]jen.Code, 1)
	templateString := entry.ReplacedTemplateValue(lang, func(key string, format dictionary.TemplateKeyFormat) string {
		formatString, argValue := g.templateFormatterCall(key, format)
		params = append(params, argValue)
		return formatString
	})
	if len(params) == 1 {
		return jen.Lit(templateString)
	} else {
		params[0] = jen.Lit(templateString)
		return jen.Qual("fmt", "Sprintf").Call(params...)
	}
}

func (g golangBuilder) templateFormatterCall(key string, templateKeys dictionary.TemplateKeyFormat) (formatString string, arg jen.Code) {
	key = code.TemplateKeyToCamelCase(key)
	switch templateKeys {
	case "int":
		return "%d", jen.Id(key)
	case "float":
		return "%f", jen.Id(key)
	case "bool":
		// TODO: Change
		return "%s", jen.Id(key)
	default:
		return "%s", jen.Id(key)
	}
}

func (g *golangBuilder) resolveEntryArgumentParam(key string, argType dictionary.TemplateKeyFormat) (callArg, paramArg *jen.Statement) {
	callArg = jen.Id(key)
	switch argType {
	case "int":
		paramArg = callArg.Clone().Int()
	case "float":
		paramArg = callArg.Clone().Float32()
	case "bool":
		paramArg = callArg.Clone().Bool()
	default:
		paramArg = callArg.Clone().String()
	}
	return
}
