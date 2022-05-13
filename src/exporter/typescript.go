package exporter

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
	"github.com/pkg/errors"
)

type TypescriptDictionaryExporter struct{}

func (t TypescriptDictionaryExporter) ExportMetadata(filePath string, metadata dictionary.Metadata) error {
	return nil
}

func (t TypescriptDictionaryExporter) ExportContent(filePath string, content dictionary.ContentRepresentation, metadata dictionary.Metadata) error {
	builder := newTypescriptContentBuilder(metadata)
	if err := builder.Run(content.ToTree()); err != nil {
		return err
	}

	file, err := os.OpenFile(path.Join(filePath, "generated/dictionary.ts"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer file.Close()
	builder.Build(metadata, file)
	return nil
}

// ExportMetadata(filePath string, metadata dictionary.Metadata) error

type TypescriptContentBuilder struct {
	dataBuilder     code.IndentedCodeBuilder
	argTypeBuilder  code.IndentedCodeBuilder
	nodeTypeBuilder code.IndentedCodeBuilder
	nodeImplBuilder code.IndentedCodeBuilder

	contentValidator dictionary.ContentValidator
	rootNode         *TypescriptContentBuilderNode
}

func (t *TypescriptContentBuilder) AddArgType(key string, value map[string]string) {
	t.argTypeBuilder.AppendLines(fmt.Sprintf("export interface %s {", key))
	t.argTypeBuilder.Indent()
	for arg, argType := range value {
		t.argTypeBuilder.AppendLines(fmt.Sprintf("%s: %s;", arg, argType))
	}
	t.argTypeBuilder.Unindent()
	t.argTypeBuilder.AppendLines("}")
}

func (t *TypescriptContentBuilder) AddNodeContent(nodeFullKey string, node *TypescriptContentBuilderNode) {
	nodeTypeInterfaceName := code.FullKeyToPascalCase(nodeFullKey) + "MDict"
	nodeImplName := nodeTypeInterfaceName + "Impl"

	t.nodeTypeBuilder.AppendLines(fmt.Sprintf("export interface %s {", nodeTypeInterfaceName))
	t.nodeTypeBuilder.Indent()
	t.nodeImplBuilder.AppendLines(fmt.Sprintf("export class %s implements %s {", nodeImplName, nodeTypeInterfaceName))
	t.nodeImplBuilder.Indent()
	t.nodeImplBuilder.AppendLines("constructor(private readonly cb: ResolverFunc) {}", "")

	for childName := range node.children {
		childTypeInterfaceName := code.FullKeyToPascalCase(childFullKey(nodeFullKey, childName)) + "MDict"
		childImplName := childTypeInterfaceName + "Impl"
		t.nodeTypeBuilder.AppendLines(fmt.Sprintf("%s: %s;", childName, childTypeInterfaceName))
		t.nodeImplBuilder.AppendLines(fmt.Sprintf("get %s() { return new %s(this.cb); }", childName, childImplName))
	}
	if len(node.children) > 0 && len(node.methods) > 0 {
		t.nodeTypeBuilder.AppendLines("")
		t.nodeImplBuilder.AppendLines("")
	}
	for methodName := range node.methods {
		fullName := node.methodFullNames[methodName]
		if node.methods[methodName] == "" {
			// No arguments
			t.nodeTypeBuilder.AppendLines(fmt.Sprintf("%s: DictionaryNFnItem;", methodName))
			t.nodeImplBuilder.AppendLines(fmt.Sprintf(`%s(language?: Language) { return this.cb("%s", undefined, language) }`, methodName, fullName))
		} else {
			t.nodeTypeBuilder.AppendLines(fmt.Sprintf("%s: DictionaryFnItem<%s>;", methodName, node.methods[methodName]))
			t.nodeImplBuilder.AppendLines(
				fmt.Sprintf(`%s(param: %s, language?: Language) { return this.cb("%s", param, language) }`,
					methodName, node.methods[methodName], fullName),
			)
		}
	}

	t.nodeTypeBuilder.Unindent()
	t.nodeTypeBuilder.AppendLines("}")
	t.nodeImplBuilder.Unindent()
	t.nodeImplBuilder.AppendLines("}")
}

func newTypescriptContentBuilder(metadata dictionary.Metadata) *TypescriptContentBuilder {
	root := newTypescriptContentBuilderNode()
	return &TypescriptContentBuilder{
		contentValidator: dictionary.NewContentValidator(metadata, dictionary.ContentValidationOptions{
			SkipLangSupportCheck: true,
		}),
		rootNode: &root,
	}
}

func (t *TypescriptContentBuilder) Run(content *dictionary.ContentNode) error {
	return t.rootNode.walk(t, content, "")
}

type TypescriptContentBuilderNode struct {
	children map[string]*TypescriptContentBuilderNode
	// Map of 'method name' -> 'argType name'. No arguments if string is zero
	methods         map[string]string
	methodFullNames map[string]string
}

func newTypescriptContentBuilderNode() TypescriptContentBuilderNode {
	return TypescriptContentBuilderNode{
		children:        map[string]*TypescriptContentBuilderNode{},
		methods:         map[string]string{},
		methodFullNames: map[string]string{},
	}
}

func (t *TypescriptContentBuilderNode) walk(builder *TypescriptContentBuilder, contentNode *dictionary.ContentNode, fullKey string) error {
	for key, child := range contentNode.Children {
		propertyName := code.FullKeyToCamelCase(key)
		childNode := newTypescriptContentBuilderNode()
		t.children[propertyName] = &childNode
		if err := childNode.walk(builder, child, childFullKey(fullKey, key)); err != nil {
			return err
		}
	}
	for key, entry := range contentNode.Entries {
		fullKey := childFullKey(fullKey, key)
		if err := builder.AddLeaf(key, fullKey, t, entry); err != nil {
			return err
		}
	}
	builder.AddNodeContent(fullKey, t)
	return nil
}

func (t *TypescriptContentBuilder) AddLeaf(key string, fullKey string, parent *TypescriptContentBuilderNode, entry dictionary.Entry) error {
	templateKeys, validateErr := t.contentValidator.Validate(entry)
	if validateErr != nil {
		return errors.Wrap(validateErr, "failed to add leaf")
	}

	interfaceName := ""
	ownArgTypeNeeded := len(templateKeys) > 0

	if ownArgTypeNeeded {
		interfaceName = t.argsInterfaceName(fullKey)
		tsArgTypes := map[string]string{}
		for k, v := range templateKeys {
			tsArgTypes[code.TemplateKeyToCamelCase(k)] = t.resolveArgumentType(v)
		}
		t.AddArgType(interfaceName, tsArgTypes)
	}
	t.AddEntry(fullKey, interfaceName, entry)

	methodName := code.FullKeyToCamelCase(key)
	parent.methods[methodName] = interfaceName
	parent.methodFullNames[methodName] = fullKey
	return nil
}

func (t *TypescriptContentBuilder) Build(metadata dictionary.Metadata, w io.Writer) {
	builder := code.IndentedCodeBuilder{}

	builder.AppendLines(
		"// Generated with donggu at "+time.Now().UTC().Format(time.RFC3339),
		"// AUTOGENERATED CODE. DO NOT EDIT.",
		"",
		`import { DictionaryFnItem, DictionaryNFnItem } from "../types";`,
		`import { Formatter } from "../util";`,
		`type ResolverFunc = (key: keyof typeof DATA, options: unknown, language?: Language) => string;`,
		"",
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

func (t *TypescriptContentBuilder) AddEntry(fullKey string, argType string, entry dictionary.Entry) {
	t.dataBuilder.AppendLines(fmt.Sprintf(`"%s": {`, fullKey))
	t.dataBuilder.Indent()
	for lang, value := range entry {
		if lang == "context" {
			continue
		}
		if argType == "" {
			t.dataBuilder.AppendLines(fmt.Sprintf("\"%s\": () => `%s`,", lang, value))
		} else {
			templateString := entry.ReplacedTemplateValue(lang, t.templateFormatterCall)
			t.dataBuilder.AppendLines(fmt.Sprintf("\"%s\": (param: %s) => `%s`,", lang, argType, templateString))
		}
	}
	t.dataBuilder.Unindent()
	t.dataBuilder.AppendLines("},")
}

func (t TypescriptContentBuilder) templateFormatterCall(key string, format dictionary.TemplateKeyFormat) string {
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

func (t TypescriptContentBuilder) argsInterfaceName(fullKey string) string {
	return fmt.Sprintf("%sArgs", code.FullKeyToPascalCase(fullKey))
}

func (t TypescriptContentBuilder) resolveArgumentType(argType dictionary.TemplateKeyFormat) string {
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

func childFullKey(key, childPart string) string {
	if key == "" {
		return childPart
	} else {
		return key + "." + childPart
	}
}
