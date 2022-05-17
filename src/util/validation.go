package util

import (
	"reflect"

	"github.com/pkg/errors"
)

type OptionMap map[string]interface{}

func SafeAccessMap[T any](obj *(map[string]interface{}), key string) (value T, err error) {
	keyValue, keyOk := (*obj)[key]
	if !keyOk {
		err = errors.Errorf("key '%s' not found", key)
		return
	}
	castValue, castOk := keyValue.(T)
	if !castOk {
		err = errors.Errorf(
			"invalid type for '%s': wanted '%s', got '%s'",
			key,
			reflect.TypeOf(*new(T)).Kind(),
			reflect.TypeOf(keyValue),
		)
		return
	} else {
		value = castValue
	}
	return
}
