package validators

import (
	"errors"
	"fmt"
	"reflect"
)

func isArray(i interface{}) bool {
	val := reflect.ValueOf(i)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		return true
	default:
		return false
	}
}

func ArrayLengthMax(i interface{}, attr map[string]interface{}) error {
	if !isArray(i) {
		return errors.New("only arrays are allowed")
	}
	if maxLen, ok := attr["max"].(float64); !ok || maxLen < 0 {
		return errors.New("attribute 'max' must be a non-negative float64")
	} else if arrLen := reflect.ValueOf(i).Len(); arrLen > int(maxLen) {
		return fmt.Errorf("array length can not be greater than %d", int(maxLen))
	}
	return nil
}

func ArrayLengthMin(i interface{}, attr map[string]interface{}) error {
	if !isArray(i) {
		return errors.New("only arrays are allowed")
	}
	if minLen, ok := attr["min"].(float64); !ok || minLen < 0 {
		return errors.New("attribute 'min' must be a non-negative float64")
	} else if arrLen := reflect.ValueOf(i).Len(); arrLen < int(minLen) {
		return fmt.Errorf("array length can not be lesser than %d", int(minLen))
	}
	return nil
}

func StringInOptions(i interface{}, attr map[string]interface{}) error {
	isString := IsString(i, attr)
	if isString != nil {
		return isString
	}
	str := i.(string)

	if _, ok := attr["options"].([]interface{}); !ok {
		return errors.New("options are required for the validator to work")
	}
	options := attr["options"].([]interface{})
	for _, op := range options {
		if o, ok := op.(string); ok {
			if o == str {
				return nil
			}
		}
	}
	return errors.New("string is out of the options")
}

func StringsExistsInOptions(i interface{}, attr map[string]interface{}) error {
	if !isArray(i) {
		return errors.New("only arrays are allowed")
	}
	STRINGS := i.([]interface{})
	for _, str := range STRINGS {
		stringDoesNotExists := StringInOptions(str, attr)
		if stringDoesNotExists != nil {
			return stringDoesNotExists
		}
	}
	return nil
}
