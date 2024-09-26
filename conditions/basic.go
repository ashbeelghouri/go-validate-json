package conditions

import (
	"github.com/ashbeelghouri/jsonschematics/structures"
)

func FieldIsProvided(_ map[string]interface{}, attr map[string]interface{}) bool {
	fields, ok := attr["schema"].(structures.Schema)
	if !ok {
		return false
	}
	toBeProvided, ok := attr["shouldBeProvided"].(string)
	if !ok {
		return false
	}
	if fields.Fields[(structures.TargetKey(toBeProvided))].Provided {
		return true
	}
	return false
}
