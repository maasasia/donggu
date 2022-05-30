package typescript

import (
	"strings"
)

func escapeTemplateStringLiteral(str string) string {
	str = strings.Replace(str, `\`, `\\`, -1)
	str = strings.Replace(str, "${", "\\${", -1)
	str = strings.Replace(str, "`", "\\`", -1)
	return str
}
