package dictionary

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Regular expression for matching numeric template options.
// Capturing groups are
// 1. flags
// 2. width
// 3. precision
const NumericOptionRegex = `([0,+]*)([1-9]+\d*)?(?:\.(\d+))?`

type TemplateKeyType string

const (
	BoolTemplateKeyType   TemplateKeyType = "bool"
	FloatTemplateKeyType  TemplateKeyType = "float"
	IntTemplateKeyType    TemplateKeyType = "int"
	StringTemplateKeyType TemplateKeyType = "string"
)

type TemplateKeyFormat struct {
	Kind   TemplateKeyType
	Option interface{}
}

type NumericTemplateFormatOption struct {
	PadCharacter   string
	WidthSet       bool
	Width          int
	PrecisionSet   bool
	Precision      int
	CommaSeparator bool
	AlwaysAddSign  bool
}

type BoolTemplateFormatOption struct {
	UseLocaleValues bool
	TrueValue       string
	FalseValue      string
}

func (t TemplateKeyFormat) Compatible(other TemplateKeyFormat) bool {
	return t.Kind == other.Kind
}

func ParseTemplateKeyFormat(kind, option string) (TemplateKeyFormat, error) {
	typedKind := TemplateKeyType(kind)
	switch typedKind {
	case BoolTemplateKeyType:
		return parseBoolFormat(option)
	case IntTemplateKeyType:
		return parseIntFormat(option)
	case FloatTemplateKeyType:
		return parseFloatFormat(option)
	case StringTemplateKeyType:
		return parseStringFormat(option)
	default:
		return TemplateKeyFormat{}, errors.Errorf("Unknown template parameter type '%s'", kind)
	}
}

func parseBoolFormat(option string) (TemplateKeyFormat, error) {
	splitOptions := strings.Split(option, ",")
	optionsIsEmpty := len(splitOptions) == 1 && strings.TrimSpace(splitOptions[0]) == ""

	if optionsIsEmpty {
		return TemplateKeyFormat{
			Kind: BoolTemplateKeyType,
			Option: BoolTemplateFormatOption{
				UseLocaleValues: true,
			},
		}, nil
	}

	if len(splitOptions) != 2 {
		return TemplateKeyFormat{}, errors.Errorf("invalid number of option values: expected 2, got %d", len(splitOptions))
	}
	return TemplateKeyFormat{
		Kind: BoolTemplateKeyType,
		Option: BoolTemplateFormatOption{
			UseLocaleValues: false,
			TrueValue:       splitOptions[0],
			FalseValue:      splitOptions[1],
		},
	}, nil
}

func parseFloatFormat(option string) (TemplateKeyFormat, error) {
	options, err := parseNumericFormat(option)
	if err != nil {
		return TemplateKeyFormat{}, err
	}
	return TemplateKeyFormat{
		Kind:   FloatTemplateKeyType,
		Option: options,
	}, nil
}

func parseIntFormat(option string) (TemplateKeyFormat, error) {
	options, err := parseNumericFormat(option)
	if err != nil {
		return TemplateKeyFormat{}, err
	}
	return TemplateKeyFormat{
		Kind:   IntTemplateKeyType,
		Option: options,
	}, nil
}

func parseStringFormat(option string) (TemplateKeyFormat, error) {
	return TemplateKeyFormat{
		Kind:   StringTemplateKeyType,
		Option: nil,
	}, nil
}

func parseNumericFormat(option string) (NumericTemplateFormatOption, error) {
	numericRegex := regexp.MustCompile(NumericOptionRegex)
	match := numericRegex.FindAllStringSubmatch(option, -1)
	if len(match) == 0 {
		return NumericTemplateFormatOption{}, errors.New("wrong format")
	}

	result := NumericTemplateFormatOption{}
	for index := range match[0][1] {
		switch match[0][1][index] {
		case '+':
			result.AlwaysAddSign = true
		case '0':
			result.PadCharacter = "0"
		case ',':
			result.CommaSeparator = true
		}
	}
	if match[0][2] != "" {
		conv, err := strconv.Atoi(match[0][2])
		if err != nil {
			return NumericTemplateFormatOption{}, errors.Wrap(err, "width is not int")
		}
		result.WidthSet = true
		result.Width = conv
	}
	if match[0][3] != "" {
		conv, err := strconv.Atoi(match[0][3])
		if err != nil {
			return NumericTemplateFormatOption{}, errors.Wrap(err, "precision is not int")
		}
		result.PrecisionSet = true
		result.Precision = conv
	}

	return result, nil
}
