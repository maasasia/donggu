package dictionary

import (
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// Regular expression for validating plural operators
const PluralOperatorRegex = `^([<>]=?)|((%|\/)[1-9]\d*)$`

type PluralDefinition struct {
	Op    string
	Value int
}

type Metadata struct {
	Version            string
	RequiredLanguages  []string
	SupportedLanguages []string
	ExporterOptions    map[string]map[string]interface{}
	Plurals            map[string][]PluralDefinition
}

func (m Metadata) SupportedLanguageSet() map[string]struct{} {
	supportedLangSet := map[string]struct{}{}
	for _, lang := range m.SupportedLanguages {
		supportedLangSet[lang] = struct{}{}
	}
	return supportedLangSet
}

func (m Metadata) RequiredLanguageSet() map[string]struct{} {
	requiredLangSet := map[string]struct{}{}
	for _, lang := range m.RequiredLanguages {
		requiredLangSet[lang] = struct{}{}
	}
	return requiredLangSet
}

func (m Metadata) ExporterOption(exporterName string) map[string]interface{} {
	if opt, ok := m.ExporterOptions[exporterName]; ok {
		return opt
	} else {
		return map[string]interface{}{}
	}
}

func (m Metadata) Validate() (err *multierror.Error) {
	supportedLangSet := map[string]struct{}{}
	requiredLangSet := map[string]struct{}{}

	if len(m.SupportedLanguages) == 0 {
		err = multierror.Append(err, errors.Errorf("supported languages is empty"))
	}
	if len(m.RequiredLanguages) == 0 {
		err = multierror.Append(err, errors.Errorf("required langauges is empty"))
	}

	for _, lang := range m.SupportedLanguages {
		if _, ok := supportedLangSet[lang]; ok {
			err = multierror.Append(err, errors.Errorf("duplicate language '%s' in SupportedLanguages", lang))
		}
		supportedLangSet[lang] = struct{}{}
		if !IsValidLanguageKey(lang) {
			err = multierror.Append(err, errors.Errorf("invalid language '%s' in SupportedLanguages", lang))
		}
	}
	for _, lang := range m.RequiredLanguages {
		if _, ok := requiredLangSet[lang]; ok {
			err = multierror.Append(err, errors.Errorf("duplicate language '%s' in RequiredLanguages", lang))
		}
		requiredLangSet[lang] = struct{}{}
		if _, ok := supportedLangSet[lang]; !ok {
			err = multierror.Append(err, errors.Errorf("language '%s' is required but not in SupportedLanguages", lang))
		}
		if !IsValidLanguageKey(lang) {
			err = multierror.Append(err, errors.Errorf("invalid language '%s' in RequiredLanguages", lang))
		}
	}
	if plError := m.validatePlurals(&supportedLangSet); plError != nil {
		err = multierror.Append(err, errors.Wrap(plError, "errors with plural definition"))
	}
	return
}

func (m Metadata) validatePlurals(languages *map[string]struct{}) (err *multierror.Error) {
	operatorRegex := regexp.MustCompile(PluralOperatorRegex)
	for lang, defs := range m.Plurals {
		if _, ok := (*languages)[lang]; !ok {
			err = multierror.Append(err, errors.Errorf("language '%s' is defined in plurals but not in SupportedLanguages", lang))
			continue
		}
		for index, def := range defs {
			if !operatorRegex.MatchString(def.Op) {
				err = multierror.Append(err, errors.Errorf("plural operator '%s' of language '%s', index [%d] is invalid", def.Op, lang, index))
			}
		}
	}
	return
}
