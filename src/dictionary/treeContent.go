package dictionary

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type ContentNode struct {
	Key      []string
	Children map[string]*ContentNode
	Entries  map[string]Entry
}

func newContent(key []string) ContentNode {
	return ContentNode{
		// TODO: Copy if key should be mutable for some reason
		Key:      key,
		Children: map[string]*ContentNode{},
		Entries:  map[string]Entry{},
	}
}

func (c *ContentNode) JoinedKey() string {
	return strings.Join(c.Key, ".")
}

func (c *ContentNode) ToTree() *ContentNode {
	return c
}

func (c *ContentNode) ToFlattened() *FlattenedContent {
	flattened := make(FlattenedContent, 0)
	c.toFlattenedWalk(&flattened)
	return &flattened
}

func (c *ContentNode) toFlattenedWalk(flattened *FlattenedContent) {
	joinedKey := c.JoinedKey()
	for key, entry := range c.Entries {
		entryKey := key
		if joinedKey != "" {
			entryKey = joinedKey + "." + entryKey
		}
		(*flattened)[entryKey] = entry
	}
	for _, child := range c.Children {
		child.toFlattenedWalk(flattened)
	}
}

func (c *ContentNode) Validate(metadata Metadata) *multierror.Error {
	validator := newContentValidator(metadata)
	return c.validateWalk(&validator)
}

func (c *ContentNode) validateWalk(validator *contentValidator) (err *multierror.Error) {
	joinedKey := c.JoinedKey()
	for key, entry := range c.Entries {
		entryKey := key
		if joinedKey != "" {
			entryKey = joinedKey + "." + entryKey
		}
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
	fmt.Printf("%s[%s]\n", indent, c.JoinedKey())
	for entry, entryContent := range c.Entries {
		fmt.Printf("%s  %s: %+v\n", indent, entry, entryContent.String())
	}
	for _, child := range c.Children {
		child.printWalk(indent + "  ")
	}
}
