package typescript

import (
	"fmt"

	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
)

type typescriptPluralBuilder struct{}

func (t typescriptPluralBuilder) Build(metadata dictionary.Metadata) code.IndentedCodeBuilder {
	builder := code.IndentedCodeBuilder{}

	builder.AppendLines("export const PLURALS = {")
	builder.Indent()

	for _, lang := range metadata.SupportedLanguages {
		if defs, ok := metadata.Plurals[lang]; ok {
			t.buildLanguage(lang, defs, &builder)
		} else {
			t.buildLanguage(lang, dictionary.DefaultPluralDefinition(), &builder)
		}
	}

	builder.Unindent()
	builder.AppendLines("};")
	return builder
}

func (t typescriptPluralBuilder) buildLanguage(lang string, defs []dictionary.PluralDefinition, builder *code.IndentedCodeBuilder) {
	builder.AppendLines(fmt.Sprintf(`"%s": (v: number) => {`, lang))
	builder.Indent()

	for index, def := range defs {
		if def.HasOperand {
			builder.AppendLines(fmt.Sprintf("if (v%s%d === %d) return %d;", def.Op, def.Operand, def.Equals, index))
		} else if def.Op == "==" {
			builder.AppendLines(fmt.Sprintf("if (v === %d) return %d;", def.Equals, index))
		} else {
			builder.AppendLines(fmt.Sprintf("if (v %s %d) return %d;", def.Op, def.Equals, index))
		}
	}
	builder.AppendLines(fmt.Sprintf("return %d;", len(defs)))

	builder.Unindent()
	builder.AppendLines("},")
}
