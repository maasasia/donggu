package dictionary

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

// contentEntryValidator is used to validate content entries with a given metadata.
// This is used to cache values derived from the metadata.
type ContentValidator struct {
	metadata          Metadata
	options           ContentValidationOptions
	supportedLangSet  map[string]struct{}
	requiredLangSet   map[string]struct{}
	templateRegex     *regexp.Regexp
	templateItemRegex *regexp.Regexp
}

type ContentValidationOptions struct {
	SkipLangSupportCheck bool
}

func NewContentValidator(m Metadata, options ContentValidationOptions) ContentValidator {
	validator := ContentValidator{
		metadata:          m,
		options:           options,
		supportedLangSet:  map[string]struct{}{},
		requiredLangSet:   map[string]struct{}{},
		templateRegex:     regexp.MustCompile("#{(.*?)}"),
		templateItemRegex: regexp.MustCompile(`#{([A-Z0-9_]+)(?:\|(string|int|float|bool))?(?:\|(.*?))?}`),
	}
	for _, lang := range m.RequiredLanguages {
		validator.requiredLangSet[lang] = struct{}{}
	}
	for _, lang := range m.SupportedLanguages {
		validator.supportedLangSet[lang] = struct{}{}
	}
	return validator
}

func (c ContentValidator) Validate(entry Entry) (templateKeys map[string]TemplateKeyFormat, err error) {
	templateKeys = map[string]TemplateKeyFormat{}
	templateKeyOwner := map[string]string{}

	if !c.options.SkipLangSupportCheck {
		for requiredLang := range c.requiredLangSet {
			if _, ok := entry[requiredLang]; !ok {
				err = errors.Errorf("'%s' is required but does not exist", requiredLang)
				return
			}
		}
	}
	for key := range entry {
		if keyErr := ValidateJoinedKey(EntryKey(key)); keyErr != nil {
			err = keyErr
			return
		}
		_, isSupportedLang := c.supportedLangSet[key]
		if !(isSupportedLang || key == "context") {
			if c.options.SkipLangSupportCheck {
				fmt.Printf("unsupported language '%s'\n", key)
			} else {
				err = errors.Errorf("language '%s' is not in supported languages", key)
				return
			}
		}
		langTemplateKeys, contentErr := entry.TemplateKeys(key)
		if contentErr != nil {
			err = errors.Wrapf(contentErr, "invalid template for '%s'", key)
			return
		}
		for templateKey, format := range langTemplateKeys {
			if existingFormat, exists := templateKeys[templateKey]; exists {
				if !format.Compatible(existingFormat) {
					err = errors.Errorf(
						"incompatible constraints in key '%s': '%s' from %s vs. '%s' from %s",
						templateKey, existingFormat, templateKeyOwner[templateKey], format, key,
					)
				}
				return
			} else {
				templateKeys[templateKey] = format
				templateKeyOwner[templateKey] = key
			}
		}
	}
	return
}

func ValidateJoinedKey(key EntryKey) error {
	for _, part := range key.Parts() {
		err := ValidateKeyPart(part)
		if err != nil {
			return errors.Wrapf(err, "invalid key '%s'", key)
		}
	}
	return nil
}

func ValidateKeySlice(key []string) error {
	for _, part := range key {
		err := ValidateKeyPart(part)
		if err != nil {
			return errors.Wrapf(err, "invalid key '%s'", strings.Join(key, "."))
		}
	}
	return nil
}

func ValidateKeyPart(keyPart string) error {
	matched, _ := regexp.MatchString("^[a-z][0-9a-z_]*$", keyPart)
	if matched {
		return nil
	} else {
		return errors.Errorf("key part '%s' is not snake_case", keyPart)
	}
}

func IsValidLanguageKey(lang string) bool {
	_, err := language.Parse(lang)
	return err == nil
}
