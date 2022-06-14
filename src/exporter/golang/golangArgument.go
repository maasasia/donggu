package golang

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
)

type golangArgumentFormatter struct{}

func (g golangArgumentFormatter) ArgumentType(key string, argType dictionary.TemplateKeyFormat) (callArg, paramArg *jen.Statement) {
	callArg = jen.Id(key)
	switch argType.Kind {
	case dictionary.IntTemplateKeyType:
		paramArg = callArg.Clone().Int()
	case dictionary.FloatTemplateKeyType:
		paramArg = callArg.Clone().Float32()
	case dictionary.BoolTemplateKeyType:
		paramArg = callArg.Clone().Bool()
	case dictionary.PluralTemplateKeyType:
		paramArg = callArg.Clone().Int()
	default:
		paramArg = callArg.Clone().String()
	}
	return
}

func (g golangArgumentFormatter) Format(key string, format dictionary.TemplateKeyFormat) (formatString string, arg jen.Code) {
	key = code.TemplateKeyToCamelCase(key)
	switch format.Kind {
	case "int":
		return g.formatNumeric(key, format) + "d", jen.Id(key)
	case "float":
		return g.formatNumeric(key, format) + "f", jen.Id(key)
	case "bool":
		return g.formatBool(key, format)
	default:
		return "%s", jen.Id(key)
	}
}

func (g golangArgumentFormatter) formatBool(key string, format dictionary.TemplateKeyFormat) (formatString string, arg jen.Code) {
	option := format.Option.(dictionary.BoolTemplateFormatOption)
	if option.UseLocaleValues {
		return "%s", jen.Id(key)
	} else {
		return "%s", jen.Id("printBooleanValue").Call(jen.Id(key), jen.Lit(option.TrueValue), jen.Lit(option.FalseValue))
	}
}

func (g golangArgumentFormatter) formatNumeric(key string, format dictionary.TemplateKeyFormat) (formatString string) {
	formatString = "%"
	option := format.Option.(dictionary.NumericTemplateFormatOption)
	if option.AlwaysAddSign {
		formatString += "+"
	}
	if option.PadCharacter == "0" {
		formatString += "0"
	}
	// TODO: Implement comma separator
	if option.WidthSet {
		formatString += fmt.Sprintf("%d", option.Width)
	}
	if option.PrecisionSet {
		formatString += "."
		formatString += fmt.Sprintf("%d", option.Precision)
	}
	return formatString
}
