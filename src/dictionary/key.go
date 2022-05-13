package dictionary

import (
	"strings"

	"github.com/maasasia/donggu/code"
)

type EntryKey string

func (e EntryKey) NewChild(childName string) EntryKey {
	if e == "" {
		return EntryKey(childName)
	} else {
		return e + "." + EntryKey(childName)
	}
}

func (e EntryKey) LastPart() string {
	return string(e[strings.LastIndex(string(e), ".")+1:])
}

func (e EntryKey) Parts() []string {
	return strings.Split(string(e), ".")
}

func (e EntryKey) PascalCase() string {
	return code.ToPascalCase(e.Parts()...)
}

func (e EntryKey) CamelCase() string {
	return code.ToCamelCase(e.Parts()...)
}

func (e EntryKey) Valid() error {
	return ValidateJoinedKey(e)
}
