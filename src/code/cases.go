package code

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func ToPascalCase(parts ...string) string {
	for i := range parts {
		partParts := strings.Split(parts[i], "_")
		for j := range partParts {
			partParts[j] = upperFirst(partParts[j])
		}
		parts[i] = strings.Join(partParts, "")
	}
	return strings.Join(parts, "")
}

func ToCamelCase(parts ...string) string {
	for i := range parts {
		partParts := strings.Split(parts[i], "_")
		for j := range partParts {
			partParts[j] = upperFirst(partParts[j])
		}
		parts[i] = strings.Join(partParts, "")
	}
	return lowerFirst(strings.Join(parts, ""))
}

func TemplateKeyToCamelCase(fullKey string) string {
	parts := strings.Split(strings.ToLower(fullKey), "_")
	for j := range parts {
		parts[j] = upperFirst(parts[j])
	}
	return lowerFirst(strings.Join(parts, ""))
}

func upperFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}
