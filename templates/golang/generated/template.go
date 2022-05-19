package generated

import (
	"fmt"
)

type ResolverFunc func(query func(lang string) bool) string

type Donggu struct {
	resolver ResolverFunc
}

func InternalNewDonggu(resolver ResolverFunc) *Donggu {
	return &Donggu{resolver: resolver}
}

func (d Donggu) resolve(key string) interface{} {
	dd, ok := formatterMappings[key]
	if !ok {
		return nil
	}
	chosenLang := d.resolver(func(lang string) bool {
		_, langExists := dd[lang]
		return langExists
	})
	if !IsValidLanguage(chosenLang) {
		panic(fmt.Errorf("language '%s' provided by resolver is invalid", chosenLang))
	}
	return dd[chosenLang]
}

func IsValidLanguage(language string) bool {
	_, ok := languages[language]
	return ok
}

func IsRequiredLanguage(language string) bool {
	lang, ok := languages[language]
	return ok && lang.bool
}

func printBooleanValue(value bool, trueValue, falseValue string) string {
	if value {
		return trueValue
	} else {
		return falseValue
	}
}
