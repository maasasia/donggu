package dictionary

import (
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
		templateItemRegex: regexp.MustCompile(`#{([A-Z0-9_]+)(?:\|(string|int|float|bool))?}`),
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
	for key, value := range entry {
		if strings.ToLower(key) != key {
			err = errors.Errorf("key '%s' should be lowercase", key)
			return
		}
		_, isSupportedLang := c.supportedLangSet[key]
		if !c.options.SkipLangSupportCheck && !(isSupportedLang || key == "context") {
			err = errors.Errorf("language '%s' is not in supported languages", key)
			return
		}
		if !isSupportedLang {
			continue
		}
		langTemplateKeys, contentErr := c.ParseFormatString(value)
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

func (c ContentValidator) ParseFormatString(format string) (map[string]TemplateKeyFormat, error) {
	templates := map[string]TemplateKeyFormat{}
	for _, template := range c.templateRegex.FindAllString(format, -1) {
		itemMatch := c.templateItemRegex.FindAllStringSubmatch(template, -1)
		if len(itemMatch) == 0 {
			return map[string]TemplateKeyFormat{}, errors.Errorf("invalid template format '%s'", template)
		}
		groups := itemMatch[0]
		if _, exists := templates[groups[1]]; exists {
			return map[string]TemplateKeyFormat{}, errors.Errorf("duplicate template key '%s'", groups[1])
		}
		if groups[2] == "" {
			templates[groups[1]] = "string"
		} else {
			templates[groups[1]] = TemplateKeyFormat(groups[2])
		}
	}
	return templates, nil
}

func ValidateJoinedKey(key string) error {
	keyParts := strings.Split(key, ".")
	for _, part := range keyParts {
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
		return errors.Errorf("invalid key part '%s'", keyPart)
	}
}

func IsValidLanguageKey(lang string) bool {
	_, err := language.Parse(lang)
	return err == nil
}
