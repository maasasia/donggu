package dictionary

import "github.com/hashicorp/go-multierror"

type ContentRepresentation interface {
	ToFlattened() *FlattenedContent
	ToTree() *ContentNode
	Validate(metadata Metadata) *multierror.Error
}

type Entry map[string]string

type TemplateKeyFormat string

func (t TemplateKeyFormat) Compatible(other TemplateKeyFormat) bool {
	return t == other
}
