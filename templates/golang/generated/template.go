package generated

type ResolverFunc func(query func(lang string) bool) string

type Donggu struct {
	Resolver ResolverFunc
}

func (d Donggu) resolve(key string) interface{} {
	dd, ok := formatterMappings[key]
	if !ok {
		return nil
	}
	return dd["ko"]
}

func IsValidLanguage(language string) bool {
	_, ok := languages[language]
	return ok
}

func IsRequiredLanguage(language string) bool {
	lang, ok := languages[language]
	return ok && lang.bool
}
