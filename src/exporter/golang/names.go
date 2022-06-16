package golang

import (
	"strings"

	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
)

func nodeStructName(key dictionary.EntryKey) string {
	return "d_" + key.PascalCase()
}

func nodeMethodChildName(child string) string {
	return code.ToPascalCase(child)
}

func entryStructName(key dictionary.EntryKey) string {
	return "d_" + key.PascalCase()
}

func entryFormatTypeName(key dictionary.EntryKey) string {
	return "d_" + key.PascalCase() + "_Fmt"
}

func entryFormatFnName(key dictionary.EntryKey, locale string) string {
	return "d_" + key.PascalCase() + "_Fmt_" + code.ToPascalCase(locale)
}

func pluralSelectorFnName(language string) string {
	return "l_plural_" + strings.ReplaceAll(language, "-", "_")
}
