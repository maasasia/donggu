package dictionary

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type Metadata struct {
	RequiredLanguages  []string
	SupportedLanguages []string
}

func (m Metadata) Validate() (err *multierror.Error) {
	supportedLangSet := map[string]struct{}{}
	requiredLangSet := map[string]struct{}{}
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
