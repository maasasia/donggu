package dictionary

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type PluralDefinition struct {
	Op         string
	Operand    int
	HasOperand bool
	Equals     int
}

func (p PluralDefinition) Valid() error {
	switch p.Op {
	case "==":
		fallthrough
	case "<":
		fallthrough
	case "<=":
		fallthrough
	case ">":
		fallthrough
	case ">=":
		return p.validateCmp()
	case "%":
		return p.validateMod()
	case "/":
		return p.validateDiv()
	default:
		return errors.Errorf("unknown operator '%s'", p.Op)
	}
}

func (p PluralDefinition) validateCmp() error {
	if p.HasOperand {
		return errors.Errorf("operator '%s' should not have operand", p.Op)
	}
	if p.Equals < 0 {
		return errors.New("value should not be negative")
	}
	return nil
}

func (p PluralDefinition) validateMod() error {
	if !p.HasOperand {
		return errors.Errorf("operator '%s' should have operand", p.Op)
	}
	if p.Operand < 0 {
		return errors.New("operand must not be negative")
	}
	if p.Operand <= p.Equals || p.Equals < 0 {
		return errors.New("value is an invalid modulo value")
	}
	return nil
}

func (p PluralDefinition) validateDiv() error {
	if !p.HasOperand {
		return errors.Errorf("operator '%s' should have operand", p.Op)
	}
	if p.Equals <= 0 || p.Operand <= 0 {
		return errors.New("value and operand must be greater than zero")
	}
	return nil
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
	for lang, defs := range m.Plurals {
		if _, ok := (*languages)[lang]; !ok {
			err = multierror.Append(err, errors.Errorf("language '%s' is defined in plurals but not in SupportedLanguages", lang))
			continue
		}
		for index, def := range defs {
			if validateErr := def.Valid(); validateErr != nil {
				err = multierror.Append(err, errors.Wrapf(validateErr, "invalid plural definition for language '%s', index [%d]", lang, index))
			}
		}
	}
	return
}

func DefaultPluralDefinition() []PluralDefinition {
	return []PluralDefinition{
		{Op: "==", Equals: 1},
	}
}
