package dictionary

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type FlattenedContent map[EntryKey]Entry

func (f *FlattenedContent) ToFlattened() *FlattenedContent {
	return f
}

func (f *FlattenedContent) ToNewFlattened() *FlattenedContent {
	copied := FlattenedContent{}
	for k, v := range *f {
		copied[k] = v
	}
	return &copied
}

// Deflatten reads itself and builds a full DictionaryContent tree.
// The receiving FlattenedDictionaryContent is expected to be valid; ie. without duplicate or invalid keys.
func (f *FlattenedContent) ToTree() *ContentNode {
	root := newContent("")
	for key, entry := range *f {
		keyParts := key.Parts()
		position := &root
		positionKey := EntryKey("")
		for i := 0; i < len(keyParts)-1; i++ {
			keyPart := keyParts[i]
			positionKey = positionKey.NewChild(keyPart)
			if position.Children[keyPart] == nil {
				newChild := newContent(positionKey)
				position.Children[keyPart] = &newChild
			}
			position = position.Children[keyPart]
		}
		position.Entries[keyParts[len(keyParts)-1]] = entry
	}
	return &root
}

func (f FlattenedContent) Validate(metadata Metadata, options ContentValidationOptions) (err *multierror.Error) {
	validator := NewContentValidator(metadata, options)
	for key, entry := range f {
		if keyErr := ValidateJoinedKey(key); keyErr != nil {
			err = multierror.Append(err, errors.Wrapf(keyErr, "invalid content '%s'", key))
		}
		if _, entryErr := validator.Validate(entry); entryErr != nil {
			err = multierror.Append(err, errors.Wrapf(entryErr, "invalid content '%s'", key))
		}
	}
	return
}
