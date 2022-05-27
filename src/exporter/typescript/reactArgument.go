package typescript

import (
	"encoding/json"
	"fmt"

	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
)

type reactArgumentFormatter struct {
	// typescriptArgumentFormatter
}

func (r reactArgumentFormatter) Format(key string, format dictionary.TemplateKeyFormat) string {
	key = code.TemplateKeyToCamelCase(key)
	switch format.Kind {
	case dictionary.FloatTemplateKeyType:
		return fmt.Sprintf("Formatter.float(param.%s, %s, options?.wrappingElement?.['%s'])", key, r.numericOptions(key, format), key)
	case dictionary.IntTemplateKeyType:
		return fmt.Sprintf("Formatter.int(param.%s, %s, options?.wrappingElement?.['%s'])", key, r.numericOptions(key, format), key)
	case dictionary.BoolTemplateKeyType:
		return r.formatBool(key, format)
	default:
		return fmt.Sprintf("Formatter.string(param.%s, options?.wrappingElement?.['%s'])", key, key)
	}
}

func (r reactArgumentFormatter) formatBool(key string, format dictionary.TemplateKeyFormat) string {
	options := format.Option.(dictionary.BoolTemplateFormatOption)
	if options.UseLocaleValues {
		return fmt.Sprintf("Formatter.bool(param.%s, options?.wrappingElement?.['%s'])", key, key)
	} else {
		return fmt.Sprintf("Formatter.string(param.%s ? `%s` : `%s`, options?.wrappingElement?.['%s'])", key, options.TrueValue, options.FalseValue, key)
	}
}

func (r reactArgumentFormatter) numericOptions(key string, format dictionary.TemplateKeyFormat) string {
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
