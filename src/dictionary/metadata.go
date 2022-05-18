package dictionary

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type Metadata struct {
	Version            string
	RequiredLanguages  []string
	SupportedLanguages []string
	ExporterOptions    map[string]map[string]interface{}
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
			err = multierror.Append(err, errors.Errorf("Language '%s' is required but not in SupportedLanguages", lang))
		}
		if !IsValidLanguageKey(lang) {
			err = multierror.Append(err, errors.Errorf("invalid language '%s' in RequiredLanguages", lang))
		}
	}
	return
}
