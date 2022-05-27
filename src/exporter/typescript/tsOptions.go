package typescript

import (
	"fmt"

	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
)

type TypescriptBuilderOptions struct{}

func (t TypescriptBuilderOptions) ArgFormatter() ArgumentFormatter {
	return typescriptArgumentFormatter{}
}

func (t TypescriptBuilderOptions) WriteEntryType(builder *code.IndentedCodeBuilder, methodName, interfaceName string, entryKey dictionary.EntryKey) {
	if interfaceName == "" {
		builder.AppendLines(fmt.Sprintf("%s: DictionaryNFnItem;", methodName))
	} else {
		builder.AppendLines(fmt.Sprintf("%s: DictionaryFnItem<%s>;", methodName, interfaceName))
	}
}

func (t TypescriptBuilderOptions) WriteEntryImpl(builder *code.IndentedCodeBuilder, methodName, interfaceName string, entryKey dictionary.EntryKey) {
	if interfaceName == "" {
		builder.AppendLines(fmt.Sprintf(`%s(language?: Language) { return this.cb("%s", undefined, language) }`, methodName, entryKey))
	} else {
		builder.AppendLines(
			fmt.Sprintf(`%s(param: %s, language?: Language) { return this.cb("%s", param, language) }`,
				methodName, interfaceName, entryKey),
		)
	}
}

func (t TypescriptBuilderOptions) WriteEntryData(builder *code.IndentedCodeBuilder, argType, language, templateString string, entry dictionary.Entry) {
	if argType == "" {
		builder.AppendLines(fmt.Sprintf("\"%s\": () => `%s`,", language, templateString))
	} else {
		templateString := entry.ReplacedTemplateValue(language, func(key string, format dictionary.TemplateKeyFormat) string {
			call := typescriptArgumentFormatter{}.Format(key, format)
			return "${" + call + "}"
		})
		builder.AppendLines(fmt.Sprintf("\"%s\": (param: %s) => `%s`,", language, argType, templateString))
	}
}
