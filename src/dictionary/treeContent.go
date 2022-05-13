package dictionary

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type ContentNode struct {
	Key      EntryKey
	Children map[string]*ContentNode
	Entries  map[string]Entry
}

func newContent(key EntryKey) ContentNode {
	return ContentNode{
		// TODO: Copy if key should be mutable for some reason
		Key:      key,
		Children: map[string]*ContentNode{},
		Entries:  map[string]Entry{},
	}
}

func (c *ContentNode) ToTree() *ContentNode {
	return c
}

func (c *ContentNode) ToNewFlattened() *FlattenedContent {
	return c.ToFlattened()
}

func (c *ContentNode) ToFlattened() *FlattenedContent {
	flattened := make(FlattenedContent, 0)
	c.toFlattenedWalk(&flattened)
	return &flattened
}

func (c *ContentNode) toFlattenedWalk(flattened *FlattenedContent) {
	for key, entry := range c.Entries {
		(*flattened)[c.Key.NewChild(key)] = entry
	}
	for _, child := range c.Children {
		child.toFlattenedWalk(flattened)
	}
}

func (c *ContentNode) Validate(metadata Metadata, options ContentValidationOptions) *multierror.Error {
	validator := NewContentValidator(metadata, options)
	return c.validateWalk(&validator)
}

func (c *ContentNode) validateWalk(validator *ContentValidator) (err *multierror.Error) {
	for key, entry := range c.Entries {
		entryKey := c.Key.NewChild(key)
		if keyErr := ValidateJoinedKey(entryKey); keyErr != nil {
			err = multierror.Append(err, errors.Wrapf(keyErr, "invalid content '%s'", key))
		}
		if _, entryErr := validator.Validate(entry); entryErr != nil {
			err = multierror.Append(err, errors.Wrapf(entryErr, "invalid content '%s'", key))
		}
	}
	for _, child := range c.Children {
		if childErr := child.validateWalk(validator); childErr != nil {
			err = multierror.Append(err, childErr)
		}
	}
	return
}

func (c *ContentNode) Print() {
	c.printWalk("")
}

func (c *ContentNode) printWalk(indent string) {
	fmt.Printf("%s[%s]\n", indent, c.Key)
	for entry, entryContent := range c.Entries {
		fmt.Printf("%s  %s: %+v\n", indent, entry, entryContent.String())
	}
	for _, child := range c.Children {
		child.printWalk(indent + "  ")
	}
}
