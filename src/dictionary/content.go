package dictionary

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
)

type ContentRepresentation interface {
	ToFlattened() *FlattenedContent
	ToTree() *ContentNode
	Validate(metadata Metadata, options ContentValidationOptions) *multierror.Error
}

type Entry map[string]string

func (e Entry) String() string {
	keys := make([]string, 0, len(e))
	for k := range e {
		keys = append(keys, k)
	}
	return fmt.Sprintf("Entry[%s]", strings.Join(keys, ", "))
}

type TemplateKeyFormat string

func (t TemplateKeyFormat) Compatible(other TemplateKeyFormat) bool {
	return t == other
}
