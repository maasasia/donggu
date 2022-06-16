package golang

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
)

type golangArgumentFormatter struct {
	metadata *dictionary.Metadata
}

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

func (g golangArgumentFormatter) Format(language, key string, format dictionary.TemplateKeyFormat) (formatString string, arg jen.Code) {
	key = code.TemplateKeyToCamelCase(key)
	switch format.Kind {
	case dictionary.IntTemplateKeyType:
		return g.formatNumeric(key, format) + "d", jen.Id(key)
	case dictionary.FloatTemplateKeyType:
		return g.formatNumeric(key, format) + "f", jen.Id(key)
	case dictionary.BoolTemplateKeyType:
		return g.formatBool(key, format)
	case dictionary.PluralTemplateKeyType:
		return "%s", g.formatPlural(language, key, format)
	default:
		return "%s", jen.Id(key)
	}
}

func (g golangArgumentFormatter) formatPlural(language, key string, format dictionary.TemplateKeyFormat) jen.Code {
	choiceStrs := format.Option.([]string)
	choices := make([]jen.Code, 0, len(choiceStrs))
	for _, choice := range choiceStrs {
		choices = append(choices, jen.Lit(choice))
	}
	return jen.Id(pluralSelectorFnName(language)).Call(jen.Id(key), jen.Index().String().Values(choices...))
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
