package typescript

import (
	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
)

type ArgumentFormatter interface {
	Format(key string, format dictionary.TemplateKeyFormat) string
}

type BuilderOptions interface {
	ArgFormatter() ArgumentFormatter
	WriteEntryType(builder *code.IndentedCodeBuilder, methodName, interfaceName string, entryKey dictionary.EntryKey)
	WriteEntryImpl(builder *code.IndentedCodeBuilder, methodName, interfaceName string, entryKey dictionary.EntryKey)
	WriteEntryData(builder *code.IndentedCodeBuilder, argType, language, templateString string, entry dictionary.Entry)
}
