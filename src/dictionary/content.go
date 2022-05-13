package dictionary

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

const (
	TemplateParenPattern  = "#{(.*?)}"
	TemplateOptionPattern = `#{([A-Z0-9_]+)(?:\|(string|int|float|bool))?}`
)

var templateParenRegex = regexp.MustCompile(TemplateParenPattern)
var templateOptionRegex = regexp.MustCompile(TemplateOptionPattern)

type ContentRepresentation interface {
	// ToFlattened returns the corresponding flattened ContentRepresentation.
	ToFlattened() *FlattenedContent
	// ToNewFlattened is same as ToFlattened, but ensures that a new copy is made.
	ToNewFlattened() *FlattenedContent
	// ToFlattened returns the corresponding tree ContentRepresentation.
	ToTree() *ContentNode
	Validate(metadata Metadata, options ContentValidationOptions) *multierror.Error
}

type TemplateKeyFormat string

func (t TemplateKeyFormat) Compatible(other TemplateKeyFormat) bool {
	return t == other
}

type Entry map[string]string

func (e Entry) TemplateKeys(key string) (map[string]TemplateKeyFormat, error) {
	templates := map[string]TemplateKeyFormat{}
	for _, template := range templateParenRegex.FindAllString(e[key], -1) {
		itemMatch := templateOptionRegex.FindAllStringSubmatch(template, -1)
		if len(itemMatch) == 0 {
			return map[string]TemplateKeyFormat{}, errors.Errorf("invalid template format '%s'", template)
		}
		groups := itemMatch[0]
		if _, exists := templates[groups[1]]; exists {
			return map[string]TemplateKeyFormat{}, errors.Errorf("duplicate template key '%s'", groups[1])
		}
		templates[groups[1]] = e.resolveKeyFormat(groups[2])
	}
	return templates, nil
}

func (e Entry) ReplacedTemplateValue(key string, replaceFn func(string, TemplateKeyFormat) string) string {
	return templateParenRegex.ReplaceAllStringFunc(e[key], func(from string) string {
		itemMatch := templateOptionRegex.FindAllStringSubmatch(from, -1)
		groups := itemMatch[0]
		return replaceFn(groups[1], e.resolveKeyFormat(groups[2]))
	})
}

func (e Entry) resolveKeyFormat(rawFormat string) TemplateKeyFormat {
	if rawFormat == "" {
		return "string"
	} else {
		return TemplateKeyFormat(rawFormat)
	}
}

func (e Entry) String() string {
	keys := make([]string, 0, len(e))
	for k := range e {
		keys = append(keys, k)
	}
	return fmt.Sprintf("Entry[%s]", strings.Join(keys, ", "))
}
