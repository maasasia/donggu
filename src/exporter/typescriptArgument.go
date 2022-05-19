package exporter

import (
	"encoding/json"
	"fmt"

	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
)

type typescriptArgumentFormatter struct{}

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
	default:
		return "string"
	}
}

func (t typescriptArgumentFormatter) Format(key string, format dictionary.TemplateKeyFormat) string {
	key = code.TemplateKeyToCamelCase(key)
	switch format.Kind {
	case dictionary.FloatTemplateKeyType:
		return fmt.Sprintf("Formatter.float(param.%s, %s)", key, t.numericOptions(key, format))
	case dictionary.IntTemplateKeyType:
		return fmt.Sprintf("Formatter.int(param.%s, %s)", key, t.numericOptions(key, format))
	case dictionary.BoolTemplateKeyType:
		return t.formatBool(key, format)
	default:
		return fmt.Sprintf("param.%s", key)
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
