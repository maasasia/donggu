package exporter

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

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
) error {
	if err := t.prepareProject(projectRoot); err != nil {
		return errors.Wrap(err, "failed to prepare project")
	}

	builder := newTypescriptContentBuilder(metadata)
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

type typescriptContentBuilder struct {
	dataBuilder     code.IndentedCodeBuilder
	argTypeBuilder  code.IndentedCodeBuilder
	nodeTypeBuilder code.IndentedCodeBuilder
	nodeImplBuilder code.IndentedCodeBuilder

	contentValidator dictionary.ContentValidator
}

func (t *typescriptContentBuilder) AddArgType(key string, value map[string]string) {
	t.argTypeBuilder.AppendLines(fmt.Sprintf("export interface %s {", key))
	t.argTypeBuilder.Indent()
	for arg, argType := range value {
		t.argTypeBuilder.AppendLines(fmt.Sprintf("%s: %s;", arg, argType))
	}
	t.argTypeBuilder.Unindent()
	t.argTypeBuilder.AppendLines("}")
}

func newTypescriptContentBuilder(metadata dictionary.Metadata) *typescriptContentBuilder {
	return &typescriptContentBuilder{
		contentValidator: dictionary.NewContentValidator(metadata, dictionary.ContentValidationOptions{
			SkipLangSupportCheck: true,
		}),
	}
}

func (t *typescriptContentBuilder) Run(content *dictionary.ContentNode) error {
	return t.walk(content, dictionary.EntryKey(""), nil)
}

func (t *typescriptContentBuilder) walk(contentNode *dictionary.ContentNode, positionKey dictionary.EntryKey, selfNameEntry dictionary.Entry) error {
	childPropertyNames := map[string]struct{}{}
	entryParamInterfaceNames := map[string]string{}
	entryFullKeys := map[string]dictionary.EntryKey{}

	entriesToSkip := map[string]struct{}{}

	for key, child := range contentNode.Children {
		propertyName := code.ToCamelCase(key)
		childPropertyNames[propertyName] = struct{}{}

		var err error
		if _, ok := contentNode.Entries[key]; ok {
			entriesToSkip[key] = struct{}{}
			err = t.walk(child, positionKey.NewChild(key), contentNode.Entries[key])
		} else {
			err = t.walk(child, positionKey.NewChild(key), nil)
		}
		if err != nil {
			return err
		}
	}
	for key, entry := range contentNode.Entries {
		if _, ok := entriesToSkip[key]; ok {
			continue
		}
		methodName, interfaceName, err := t.addEntry(entry, positionKey.NewChild(key), false)
		if err != nil {
			return err
		}
		entryParamInterfaceNames[methodName] = interfaceName
		entryFullKeys[methodName] = positionKey.NewChild(key)
	}
	if selfNameEntry != nil {
		methodName, interfaceName, err := t.addEntry(selfNameEntry, positionKey.NewChild("$"), true)
		if err != nil {
			return err
		}
		entryParamInterfaceNames[methodName] = interfaceName
		entryFullKeys[methodName] = positionKey.NewChild("$")
	}

	t.writeNodeToBuilder(positionKey, &childPropertyNames, &entryParamInterfaceNames, &entryFullKeys)
	return nil
}

func (t *typescriptContentBuilder) addEntry(entry dictionary.Entry, entryKey dictionary.EntryKey, isSelfEntry bool) (methodName string, interfaceName string, err error) {
	templateKeys, validateErr := t.contentValidator.Validate(entry)
	if validateErr != nil {
		err = errors.Wrap(validateErr, "failed to add leaf")
		return
	}

	methodName = "_"
	if !isSelfEntry {
		methodName = code.ToCamelCase(entryKey.LastPart())
	}
	interfaceName = ""
	ownArgTypeNeeded := len(templateKeys) > 0

	if ownArgTypeNeeded {
		interfaceName = t.argsInterfaceName(entryKey)
		tsArgTypes := map[string]string{}
		for k, v := range templateKeys {
			tsArgTypes[code.TemplateKeyToCamelCase(k)] = t.resolveArgumentType(v)
		}
		t.AddArgType(interfaceName, tsArgTypes)
	}
	t.writeEntryDataToBuilder(entryKey, interfaceName, entry)

	return
}

func (t *typescriptContentBuilder) writeNodeToBuilder(
	parentKey dictionary.EntryKey,
	childPropertyNames *map[string]struct{},
	entryParamInterfaceNames *map[string]string,
	entryFullKeys *map[string]dictionary.EntryKey) {

	nodeTypeInterfaceName := t.nodeInterfaceName(parentKey)
	nodeImplName := t.nodeImplName(parentKey)

	t.nodeTypeBuilder.AppendLines(fmt.Sprintf("export interface %s {", nodeTypeInterfaceName))
	t.nodeTypeBuilder.Indent()
	t.nodeImplBuilder.AppendLines(fmt.Sprintf("export class %s implements %s {", nodeImplName, nodeTypeInterfaceName))
	t.nodeImplBuilder.Indent()
	t.nodeImplBuilder.AppendLines("constructor(private readonly cb: ResolverFunc) {}", "")

	for childName := range *childPropertyNames {
		childKey := parentKey.NewChild(childName)
		childTypeInterfaceName := t.nodeInterfaceName(childKey)
		childImplName := t.nodeImplName(childKey)
		t.nodeTypeBuilder.AppendLines(fmt.Sprintf("%s: %s;", childName, childTypeInterfaceName))
		t.nodeImplBuilder.AppendLines(fmt.Sprintf("get %s() { return new %s(this.cb); }", childName, childImplName))
	}
	if len(*childPropertyNames) > 0 && len(*entryFullKeys) > 0 {
		t.nodeTypeBuilder.AppendLines("")
		t.nodeImplBuilder.AppendLines("")
	}
	for methodName := range *entryFullKeys {
		entryKey := (*entryFullKeys)[methodName]
		interfaceName := (*entryParamInterfaceNames)[methodName]
		if interfaceName == "" {
			// No arguments
			t.nodeTypeBuilder.AppendLines(fmt.Sprintf("%s: DictionaryNFnItem;", methodName))
			t.nodeImplBuilder.AppendLines(fmt.Sprintf(`%s(language?: Language) { return this.cb("%s", undefined, language) }`, methodName, entryKey))
		} else {
			t.nodeTypeBuilder.AppendLines(fmt.Sprintf("%s: DictionaryFnItem<%s>;", methodName, interfaceName))
			t.nodeImplBuilder.AppendLines(
				fmt.Sprintf(`%s(param: %s, language?: Language) { return this.cb("%s", param, language) }`,
					methodName, interfaceName, entryKey),
			)
		}
	}

	t.nodeTypeBuilder.Unindent()
	t.nodeTypeBuilder.AppendLines("}")
	t.nodeImplBuilder.Unindent()
	t.nodeImplBuilder.AppendLines("}")
}

func (t *typescriptContentBuilder) writeEntryDataToBuilder(fullKey dictionary.EntryKey, argType string, entry dictionary.Entry) {
	t.dataBuilder.AppendLines(fmt.Sprintf(`"%s": {`, fullKey))
	t.dataBuilder.Indent()
	for lang, value := range entry {
		if lang == "context" {
			continue
		}
		if argType == "" {
			t.dataBuilder.AppendLines(fmt.Sprintf("\"%s\": () => `%s`,", lang, value))
		} else {
			templateString := entry.ReplacedTemplateValue(lang, func(key string, format dictionary.TemplateKeyFormat) string {
				return "${" + t.templateFormatterCall(key, format) + "}"
			})
			t.dataBuilder.AppendLines(fmt.Sprintf("\"%s\": (param: %s) => `%s`,", lang, argType, templateString))
		}
	}
	t.dataBuilder.Unindent()
	t.dataBuilder.AppendLines("},")
}

func (t typescriptContentBuilder) templateFormatterCall(key string, format dictionary.TemplateKeyFormat) string {
	key = code.TemplateKeyToCamelCase(key)
	switch format {
	case "int":
		return fmt.Sprintf("Formatter.int(param.%s)", key)
	case "float":
		return fmt.Sprintf("Formatter.float(param.%s)", key)
	case "bool":
		return fmt.Sprintf("Formatter.bool(param.%s)", key)
	default:
		return fmt.Sprintf("param.%s", key)
	}
}

func (t *typescriptContentBuilder) Build(metadata dictionary.Metadata, w io.Writer) {
	builder := code.IndentedCodeBuilder{}

	builder.AppendLines(
		"// Generated with donggu at "+time.Now().UTC().Format(time.RFC3339),
		"// AUTOGENERATED CODE. DO NOT EDIT.",
		"",
		`import { DictionaryFnItem, DictionaryNFnItem } from "../types";`,
		`import { Formatter } from "../util";`,
		`type ResolverFunc = (key: keyof typeof DATA, options: unknown, language?: Language) => string;`,
		"",
		fmt.Sprintf("export const Version = '%s';", metadata.Version),
		fmt.Sprintf("export type RequiredLanguage = '%s';", strings.Join(metadata.RequiredLanguages, "' | '")),
		fmt.Sprintf("export type Language = '%s';", strings.Join(metadata.SupportedLanguages, "' | '")),
		"",
	)

	builder.AppendLines("export const DATA = {")
	builder.IndentedBlock(t.dataBuilder)
	builder.AppendLines("};", "")
	builder.AppendBlock(t.argTypeBuilder)
	builder.AppendLines("")
	builder.AppendBlock(t.nodeTypeBuilder)
	builder.AppendLines("")
	builder.AppendBlock(t.nodeImplBuilder)
	builder.AppendLines("")

	builder.Build(w)
}

func (t typescriptContentBuilder) argsInterfaceName(fullKey dictionary.EntryKey) string {
	return fullKey.PascalCase() + "_Args"
}

func (t typescriptContentBuilder) nodeInterfaceName(key dictionary.EntryKey) string {
	return key.PascalCase() + "_MDict"
}

func (t typescriptContentBuilder) nodeImplName(key dictionary.EntryKey) string {
	return t.nodeInterfaceName(key) + "_Impl"
}

func (t typescriptContentBuilder) resolveArgumentType(argType dictionary.TemplateKeyFormat) string {
	switch argType {
	case "int":
		fallthrough
	case "float":
		return "number"
	case "bool":
		return "boolean"
	default:
		return "string"
	}
}
