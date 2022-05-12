package dictionary

func MergeContent(from, to ContentRepresentation) ContentRepresentation {
	flatFrom, flatTo := from.ToNewFlattened(), from.ToNewFlattened()
	for fromLang, fromEntry := range *flatFrom {
		if _, ok := (*flatTo)[fromLang]; ok {
			for lang, langContent := range fromEntry {
				(*flatTo)[fromLang][lang] = langContent
			}
		} else {
			(*flatTo)[fromLang] = fromEntry
		}
	}
	return flatTo
}
