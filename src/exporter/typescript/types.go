package typescript

import (
	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
	"github.com/maasasia/donggu/util"
)

type ArgumentFormatter interface {
	Format(language, key string, format dictionary.TemplateKeyFormat) string
}

type BuilderOptions interface {
	SetShortener(shortener util.Shortener)
	ArgFormatter() ArgumentFormatter
	WriteHeader(builder *code.IndentedCodeBuilder)
	WriteEntryType(builder *code.IndentedCodeBuilder, methodName, interfaceName string, entryKey dictionary.EntryKey)
	WriteEntryImpl(builder *code.IndentedCodeBuilder, methodName, interfaceName string, entryKey dictionary.EntryKey)
	WriteEntryData(builder *code.IndentedCodeBuilder, argType, language, templateString string, entry dictionary.Entry) error
}
