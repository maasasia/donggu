package dictionary

type DiffFlag string

const (
	CreateFlag DiffFlag = "CREATE"
	DeleteFlag DiffFlag = "DELETE"
	ChangeFlag DiffFlag = "CHANGE"
)

type ContentDifference struct {
	CreatedKeys       []EntryKey
	DeletedKeys       []EntryKey
	Changes           map[EntryKey]map[string]DiffFlag
	SourceKeyCount    int
	DestKeyCount      int
	UnchangedKeyCount int
}

func DiffContents(from, to ContentRepresentation) ContentDifference {
	flatFrom, flatTo := from.ToFlattened(), to.ToFlattened()
	diff := ContentDifference{
		CreatedKeys:       []EntryKey{},
		DeletedKeys:       []EntryKey{},
		Changes:           map[EntryKey]map[string]DiffFlag{},
		SourceKeyCount:    len(*flatFrom),
		DestKeyCount:      len(*flatTo),
		UnchangedKeyCount: 0,
	}

	for key := range *flatFrom {
		if _, ok := (*flatTo)[key]; !ok {
			diff.DeletedKeys = append(diff.DeletedKeys, key)
		}
	}
	for key := range *flatTo {
		if _, ok := (*flatFrom)[key]; !ok {
			diff.CreatedKeys = append(diff.CreatedKeys, key)
		} else {
			if len((*flatTo)[key]) == 0 && len((*flatFrom)[key]) == 0 {
				continue
			}
			if len((*flatTo)[key]) == 0 {
				diff.DeletedKeys = append(diff.DeletedKeys, key)
				continue
			}
			if len((*flatFrom)[key]) == 0 {
				diff.CreatedKeys = append(diff.CreatedKeys, key)
				continue
			}

			diff.Changes[key] = map[string]DiffFlag{}
			for lang := range (*flatFrom)[key] {
				if _, ok := (*flatTo)[key][lang]; !ok {
					diff.Changes[key][lang] = DeleteFlag
				}
			}
			for lang, langValue := range (*flatTo)[key] {
				if langFromValue, ok := (*flatFrom)[key][lang]; !ok {
					diff.Changes[key][lang] = CreateFlag
				} else if langValue != langFromValue {
					diff.Changes[key][lang] = ChangeFlag
				}
			}
			if len(diff.Changes[key]) == 0 {
				delete(diff.Changes, key)
			}
		}
	}
	diff.UnchangedKeyCount = len(*flatTo) - len(diff.Changes) - len(diff.CreatedKeys)
	return diff
}
