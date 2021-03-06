package golang

import (
	"os"
	"path"
	"time"

	"github.com/dave/jennifer/jen"
	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
	"github.com/pkg/errors"
)

type GolangBuilder struct {
	builder          *golangCodeBuilder
	metadata         *dictionary.Metadata
	contentValidator dictionary.ContentValidator
}

func NewGolangBuilder(metadata dictionary.Metadata) *GolangBuilder {
	return &GolangBuilder{
		builder:  newGolangCodeBuilder(),
		metadata: &metadata,
		contentValidator: dictionary.NewContentValidator(metadata, dictionary.ContentValidationOptions{
			SkipLangSupportCheck: true,
		}),
	}
}

func (g *GolangBuilder) Build(metadata dictionary.Metadata, projectRoot string) error {
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

func (g *GolangBuilder) openFile(projectRoot, filename string) (*os.File, error) {
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

func (g *GolangBuilder) Run(content *dictionary.ContentNode) error {
	return g.walk(content, dictionary.EntryKey(""), nil, 0)
}

func (g *GolangBuilder) walk(
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

func (g *GolangBuilder) addEntry(entry dictionary.Entry, entryKey dictionary.EntryKey, isSelfEntry bool) (err error) {
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
		formatterValue, formatErr := g.buildFormatterReturnValue(entry, lang, templateKeys)
		if err != nil {
			err = errors.Wrap(formatErr, "failed to build formatter value")
			return
		}
		g.builder.writeEntryImpl(entryKey, lang, paramArgs, formatterValue)
	}
	return
}

func (g *GolangBuilder) buildEntryArgumentBlock(argTypes map[string]dictionary.TemplateKeyFormat) (callArgs, paramArgs []jen.Code) {
	callArgs = make([]jen.Code, 0, len(argTypes))
	paramArgs = make([]jen.Code, 0, len(argTypes))
	for k, v := range argTypes {
		callArg, paramArg := golangArgumentFormatter{metadata: g.metadata}.ArgumentType(code.TemplateKeyToCamelCase(k), v)
		callArgs = append(callArgs, callArg)
		paramArgs = append(paramArgs, paramArg)
	}
	return
}
func (g *GolangBuilder) writeNodeToBuilder(
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

func (g *GolangBuilder) buildFormatterReturnValue(
	entry dictionary.Entry,
	lang string,
	argTypes map[string]dictionary.TemplateKeyFormat,
) (*jen.Statement, error) {
	params := make([]jen.Code, 1)
	templateString, err := entry.ReplacedTemplateValue(lang, func(key string, format dictionary.TemplateKeyFormat) (string, error) {
		formatString, argValue := golangArgumentFormatter{}.Format(lang, key, format)
		params = append(params, argValue)
		return formatString, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse template parameter")
	}
	if len(params) == 1 {
		return jen.Lit(templateString), nil
	} else {
		params[0] = jen.Lit(templateString)
		return jen.Qual("fmt", "Sprintf").Call(params...), nil
	}
}
