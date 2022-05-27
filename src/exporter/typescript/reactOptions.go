package typescript

import (
	"fmt"
	"time"

	"github.com/maasasia/donggu/code"
	"github.com/maasasia/donggu/dictionary"
	"github.com/maasasia/donggu/util"
)

type ReactBuilderOptions struct{ shortener util.Shortener }

func (r *ReactBuilderOptions) SetShortener(shortener util.Shortener) {
	r.shortener = shortener
}

func (r ReactBuilderOptions) ArgFormatter() ArgumentFormatter {
	return reactArgumentFormatter{}
}

func (t ReactBuilderOptions) WriteHeader(builder *code.IndentedCodeBuilder) {
	builder.AppendLines(
		"// Generated with donggu at "+time.Now().UTC().Format(time.RFC3339),
		"// AUTOGENERATED CODE. DO NOT EDIT.",
		"",
		`import React from "react";`,
		"",
		`import { DictionaryFnItem, DictionaryNFnItem, EntryOptions } from "../types";`,
		`import { Formatter, replaceLineBreak as rlb } from "../util";`,
		`type ResolverFunc = (key: keyof typeof DATA, options: unknown, language?: Language) => string;`,
		"",
	)
}

func (t ReactBuilderOptions) WriteEntryType(builder *code.IndentedCodeBuilder, methodName, interfaceName string, entryKey dictionary.EntryKey) {
	if interfaceName == "" {
		builder.AppendLines(fmt.Sprintf("%s: DictionaryNFnItem;", methodName))
	} else {
		builder.AppendLines(fmt.Sprintf("%s: DictionaryFnItem<%s>;", methodName, interfaceName))
	}
}

func (t ReactBuilderOptions) WriteEntryImpl(builder *code.IndentedCodeBuilder, methodName, interfaceName string, entryKey dictionary.EntryKey) {
	if interfaceName == "" {
		builder.AppendLines(fmt.Sprintf(`%s(options?: EntryOptions) { return this.cb("%s", options) }`, methodName, t.shortener.Shorten(string(entryKey))))
	} else {
		builder.AppendLines(
			fmt.Sprintf(`%s(param: %s, options?: EntryOptions<%s>) { return this.cb("%s", options) }`,
				methodName, interfaceName, interfaceName, t.shortener.Shorten(string(entryKey))),
		)
	}
}

func (t ReactBuilderOptions) WriteEntryData(builder *code.IndentedCodeBuilder, argType, language, templateString string, entry dictionary.Entry) {
	if argType == "" {
		builder.AppendLines(fmt.Sprintf("\"%s\": (options?: EntryOptions) => <>{rlb(`%s`,options?.lineBreakElement)}</>,", language, templateString))
	} else {
		templateString := entry.ReplacedTemplateValue(language, func(key string, format dictionary.TemplateKeyFormat) string {
			call := reactArgumentFormatter{}.Format(key, format)
			return "`,options?.lineBreakElement)}{" + call + "}{rlb(`"
		})
		builder.AppendLines(fmt.Sprintf("\"%s\": (param: %s, options?: EntryOptions<%s>) => <>{rlb(`%s`,options?.lineBreakElement)}</>,", language, argType, argType, templateString))
	}
}
