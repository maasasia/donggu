package typescript

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
)

type typescriptArgumentFormatter struct {
	metadata *dictionary.Metadata
}

type typescriptNumericFormatJsonMarshal struct {
	PadCharacter   *string `json:"padCharacter"`
	Width          *int    `json:"width"`
	Precision      *int    `json:"precision"`
	CommaSeparator bool    `json:"comma"`
	AlwaysAddSign  bool    `json:"alwaysSign"`
}

func (t typescriptArgumentFormatter) ArgumentType(format dictionary.TemplateKeyFormat) string {
	switch format.Kind {
	case dictionary.FloatTemplateKeyType:
		fallthrough
	case dictionary.IntTemplateKeyType:
		return "number"
	case dictionary.BoolTemplateKeyType:
		return "boolean"
	case dictionary.PluralTemplateKeyType:
		return "number"
	default:
		return "string"
	}
}

func (t typescriptArgumentFormatter) Format(language, key string, format dictionary.TemplateKeyFormat) (string, error) {
	key = code.TemplateKeyToCamelCase(key)
	switch format.Kind {
	case dictionary.FloatTemplateKeyType:
		return fmt.Sprintf("Formatter.float(param.%s, %s)", key, t.numericOptions(key, format)), nil
	case dictionary.IntTemplateKeyType:
		return fmt.Sprintf("Formatter.int(param.%s, %s)", key, t.numericOptions(key, format)), nil
	case dictionary.BoolTemplateKeyType:
		return t.formatBool(key, format), nil
	case dictionary.PluralTemplateKeyType:
		return t.formatPlural(language, key, format)
	default:
		return fmt.Sprintf("param.%s", key), nil
	}
}

func (t typescriptArgumentFormatter) formatBool(key string, format dictionary.TemplateKeyFormat) string {
	options := format.Option.(dictionary.BoolTemplateFormatOption)
	if options.UseLocaleValues {
		return fmt.Sprintf("Formatter.bool(param.%s)", key)
	} else {
		return fmt.Sprintf("param.%s ? `%s` : `%s`", key, options.TrueValue, options.FalseValue)
	}
}

func (t typescriptArgumentFormatter) formatPlural(language, key string, format dictionary.TemplateKeyFormat) (string, error) {
	options := format.Option.([]string)
	optArray := strings.Join(options, `","`)
	if !checkPluralOptionLength(format, language, t.metadata) {
		return "", errors.New("plural option length does not match")
	}
	return fmt.Sprintf(`Formatter.plural(param.%s, "%s", ["%s"])`, key, language, optArray), nil
}

func (t typescriptArgumentFormatter) numericOptions(key string, format dictionary.TemplateKeyFormat) string {
	options := format.Option.(dictionary.NumericTemplateFormatOption)
	if options.IsZero() {
		return "null"
	}
	optionMarshal := typescriptNumericFormatJsonMarshal{
		PadCharacter:   nil,
		CommaSeparator: options.CommaSeparator,
		AlwaysAddSign:  options.AlwaysAddSign,
		Width:          nil,
		Precision:      nil,
	}
	if options.WidthSet || options.PadCharacter == "0" {
		optionMarshal.PadCharacter = new(string)
		if options.PadCharacter == "0" {
			*optionMarshal.PadCharacter = "0"
		} else {
			*optionMarshal.PadCharacter = " "
		}
	}
	if options.WidthSet {
		optionMarshal.Width = new(int)
		*optionMarshal.Width = options.Width
	}
	if options.PrecisionSet {
		optionMarshal.Precision = new(int)
		*optionMarshal.Precision = options.Precision
	}
	optString, _ := json.Marshal(optionMarshal)
	return string(optString)
}
