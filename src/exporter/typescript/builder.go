package typescript

import (
	"fmt"
	"io"
	"strings"

	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
	"github.com/pkg/errors"
)

type typescriptBuilder struct {
	options          BuilderOptions
	contentValidator dictionary.ContentValidator

	dataBuilder     code.IndentedCodeBuilder
	argTypeBuilder  code.IndentedCodeBuilder
	nodeTypeBuilder code.IndentedCodeBuilder
	nodeImplBuilder code.IndentedCodeBuilder
}

func NewTypescriptBuilder(metadata dictionary.Metadata, options BuilderOptions) *typescriptBuilder {
	return &typescriptBuilder{
		contentValidator: dictionary.NewContentValidator(metadata, dictionary.ContentValidationOptions{
			SkipLangSupportCheck: true,
		}),
		options: options,
	}
}

func (t *typescriptBuilder) AddArgType(key string, value map[string]string) {
	t.argTypeBuilder.AppendLines(fmt.Sprintf("export interface %s {", key))
	t.argTypeBuilder.Indent()
	for arg, argType := range value {
		t.argTypeBuilder.AppendLines(fmt.Sprintf("%s: %s;", arg, argType))
	}
	t.argTypeBuilder.Unindent()
	t.argTypeBuilder.AppendLines("}")
}

func (t *typescriptBuilder) Run(content *dictionary.ContentNode) error {
	return t.walk(content, dictionary.EntryKey(""), nil)
}

func (t *typescriptBuilder) walk(contentNode *dictionary.ContentNode, positionKey dictionary.EntryKey, selfNameEntry dictionary.Entry) error {
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

func (t *typescriptBuilder) addEntry(entry dictionary.Entry, entryKey dictionary.EntryKey, isSelfEntry bool) (methodName string, interfaceName string, err error) {
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
			argType := typescriptArgumentFormatter{}.ArgumentType(v)
			tsArgTypes[code.TemplateKeyToCamelCase(k)] = argType
		}
		t.AddArgType(interfaceName, tsArgTypes)
	}
	t.writeEntryDataToBuilder(entryKey, interfaceName, entry)

	return
}

func (t *typescriptBuilder) writeNodeToBuilder(
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

		t.options.WriteEntryType(&t.nodeTypeBuilder, methodName, interfaceName, entryKey)
		t.options.WriteEntryImpl(&t.nodeImplBuilder, methodName, interfaceName, entryKey)
	}

	t.nodeTypeBuilder.Unindent()
	t.nodeTypeBuilder.AppendLines("}")
	t.nodeImplBuilder.Unindent()
	t.nodeImplBuilder.AppendLines("}")
}

func (t *typescriptBuilder) writeEntryDataToBuilder(fullKey dictionary.EntryKey, argType string, entry dictionary.Entry) {
	t.dataBuilder.AppendLines(fmt.Sprintf(`"%s": {`, fullKey))
	t.dataBuilder.Indent()
	for lang, value := range entry {
		if lang == "context" {
			continue
		}
		t.options.WriteEntryData(&t.dataBuilder, argType, lang, value, entry)
	}
	t.dataBuilder.Unindent()
	t.dataBuilder.AppendLines("},")
}

func (t *typescriptBuilder) Build(metadata dictionary.Metadata, w io.Writer) {
	builder := code.IndentedCodeBuilder{}
	t.options.WriteHeader(&builder)
	builder.AppendLines(
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

func (t typescriptBuilder) argsInterfaceName(fullKey dictionary.EntryKey) string {
	return fullKey.PascalCase() + "_Args"
}

func (t typescriptBuilder) nodeInterfaceName(key dictionary.EntryKey) string {
	return key.PascalCase() + "_MDict"
}

func (t typescriptBuilder) nodeImplName(key dictionary.EntryKey) string {
	return t.nodeInterfaceName(key) + "_Impl"
}