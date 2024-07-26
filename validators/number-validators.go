package validators

import (
	"errors"
	"fmt"
	"reflect"
)

var NumberTypes = map[string][]string{
	"integer": {
		"int",
		"int32",
		"int64",
	},
	"float": {
		"float",
		"float32",
		"float64",
	},
}

func convertToFloat64(i interface{}) *float64 {
	var result float64
	switch v := i.(type) {
	case int:
		result = float64(v)
	case int8:
		result = float64(v)
	case int16:
		result = float64(v)
	case int32:
		result = float64(v)
	case int64:
		result = float64(v)
	case uint:
		result = float64(v)
	case uint8:
		result = float64(v)
	case uint16:
		result = float64(v)
	case uint32:
		result = float64(v)
	case uint64:
		result = float64(v)
	case float32:
		result = float64(v)
	case float64:
		result = float64(v)
	default:
		return nil
	}
	return &result
}

func IsInteger(i interface{}, _ map[string]interface{}) error {
	typeOfInterface := reflect.TypeOf(i).String()

	for _, t := range NumberTypes["integer"] {
		if t == typeOfInterface {
			return nil
		}
	}

	return errors.New("value is not an integer")
}

func IsFloat(i interface{}, _ map[string]interface{}) error {
	typeOfInterface := reflect.TypeOf(i).String()

	for _, t := range NumberTypes["float"] {
		if t == typeOfInterface {
			return nil
		}
	}

	return errors.New("value is not a floating number")
}

func IsNumber(i interface{}, attr map[string]interface{}) error {
	if err := IsInteger(i, attr); err == nil {
		return nil
	}
	if err := IsFloat(i, attr); err == nil {
		return nil
	}
	return errors.New("value is neither integer not floating number")
}

func MaxAllowed(i interface{}, attributes map[string]interface{}) error {
	number := convertToFloat64(i)
	if number == nil {
		return errors.New(fmt.Sprintf("%v is not a number", i))
	}
	if _, ok := attributes["max"]; !ok {
		return errors.New("max attribute is required")
	}
	_max := convertToFloat64(attributes["max"])
	if _max == nil {
		return errors.New("max attribute should be a number")
	}
	if *number > *_max {
		return errors.New(fmt.Sprintf("%d is greater than %d", *number, *_max))
	}
	return nil
}

func MinAllowed(i interface{}, attributes map[string]interface{}) error {
	number := convertToFloat64(i)
	if number == nil {
		return errors.New(fmt.Sprintf("%v is not a number", i))
	}
	if _, ok := attributes["min"]; !ok {
		return errors.New("min attribute is required")
	}
	_max := convertToFloat64(attributes["min"])
	if _max == nil {
		return errors.New("min attribute should be a number")
	}
	if *number < *_max {
		return errors.New(fmt.Sprintf("%d is lesser than %d", *number, *_max))
	}
	return nil
}

func InBetween(i interface{}, attributes map[string]interface{}) error {
	if err := MinAllowed(i, attributes); err != nil {
		return err
	}
	if err := MaxAllowed(i, attributes); err != nil {
		return err
	}
	return nil
}

func IsGreaterThanZero(i interface{}, _ map[string]interface{}) error {
	return MinAllowed(i, map[string]interface{}{
		"min": 0,
	})
}

func IsLesserThanZero(i interface{}, _ map[string]interface{}) error {
	return MinAllowed(i, map[string]interface{}{
		"max": 0,
	})
}
