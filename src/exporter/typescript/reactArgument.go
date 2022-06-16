package typescript

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
)

type reactArgumentFormatter struct {
	metadata *dictionary.Metadata
}

func (r reactArgumentFormatter) Format(language, key string, format dictionary.TemplateKeyFormat) (string, error) {
	key = code.TemplateKeyToCamelCase(key)
	switch format.Kind {
	case dictionary.FloatTemplateKeyType:
		ret := fmt.Sprintf("Formatter.float(param.%s, %s, options?.wrappingElement?.['%s'])", key, r.numericOptions(key, format), key)
		return ret, nil
	case dictionary.IntTemplateKeyType:
		ret := fmt.Sprintf("Formatter.int(param.%s, %s, options?.wrappingElement?.['%s'])", key, r.numericOptions(key, format), key)
		return ret, nil
	case dictionary.BoolTemplateKeyType:
		return r.formatBool(key, format), nil
	case dictionary.PluralTemplateKeyType:
		return r.formatPlural(language, key, format)
	default:
		return fmt.Sprintf("Formatter.string(param.%s, options?.wrappingElement?.['%s'])", key, key), nil
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

func (r reactArgumentFormatter) formatPlural(language, key string, format dictionary.TemplateKeyFormat) (string, error) {
	options := format.Option.([]string)
	optArray := strings.Join(options, `","`)
	if !checkPluralOptionLength(format, language, r.metadata) {
		return "", errors.New("plural option length does not match")
	}
	return fmt.Sprintf(`Formatter.plural(param.%s, "%s", ["%s"])`, key, language, optArray), nil
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
