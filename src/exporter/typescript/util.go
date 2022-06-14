package typescript

import (
	"strings"

	"github.com/maasasia/donggu/dictionary"
)

func escapeTemplateStringLiteral(str string) string {
	str = strings.Replace(str, `\`, `\\`, -1)
	str = strings.Replace(str, "${", "\\${", -1)
	str = strings.Replace(str, "`", "\\`", -1)
	return str
}

func checkPluralOptionLength(format dictionary.TemplateKeyFormat, language string, metadata *dictionary.Metadata) bool {
	optionLength := len(format.Option.([]string))
	if defs, ok := metadata.Plurals[language]; ok {
		return optionLength == len(defs)+1
	} else {
		return optionLength == 2
	}
}
