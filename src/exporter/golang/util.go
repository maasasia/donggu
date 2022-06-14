package golang

import "github.com/maasasia/donggu/dictionary"

func checkPluralOptionLength(format dictionary.TemplateKeyFormat, language string, metadata *dictionary.Metadata) bool {
	optionLength := len(format.Option.([]string))
	if defs, ok := metadata.Plurals[language]; ok {
		return optionLength == len(defs)+1
	} else {
		return optionLength == 2
	}
}
