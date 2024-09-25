package v0

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ashbeelghouri/jsonschematics/conditions"
	"github.com/ashbeelghouri/jsonschematics/errorHandler"
	"github.com/ashbeelghouri/jsonschematics/operators"
	"github.com/ashbeelghouri/jsonschematics/utils"
	"github.com/ashbeelghouri/jsonschematics/validators"
	"log"
	"os"
	"strings"
	"sync"
)

type TargetKey string

type Schematics struct {
	Schema     Schema
	Validators validators.Validators
	Operators  operators.Operators
	Conditions conditions.Conditions
	Separator  string
	ArrayIdKey string
	Locale     string
	DB         map[string]interface{}
	FlatData   map[string]interface{}
	UnFlatData map[string]interface{}
	Logging    utils.Logger
}

// add this DB to the attributes as SCHEMA_GLOBAL_DB

type Schema struct {
	Version string                 `json:"version"`
	Fields  map[TargetKey]Field    `json:"fields"`
	DB      map[string]interface{} `json:"DB"`
}

type Field struct {
	DependsOn             []string               `json:"depends_on"`
	Target                string                 `json:"target"`
	DisplayName           string                 `json:"display_name"`
	Name                  string                 `json:"name"`
	Type                  string                 `json:"type"`
	IsRequired            bool                   `json:"required"`
	AddToDB               bool                   `json:"add_to_db"`
	Description           string                 `json:"description"`
	Validators            map[string]Constant    `json:"validators"`
	Operators             map[string]Constant    `json:"operators"`
	L10n                  map[string]interface{} `json:"l10n"`
	AdditionalInformation map[string]interface{} `json:"additional_information"`
	Conditions            map[string]Condition   `json:"conditions"`
	Tags                  []string               `json:"tags"`
	Value                 map[string]interface{} `json:"value"`
	Provided              bool
	Status                string
	Errors                errorHandler.Errors
	logging               utils.Logger
}

func (f *Field) AsMap() *map[string]interface{} {
	var _map map[string]interface{}
	_bytes, err := json.Marshal(f)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(_bytes, &_map)
	if err != nil {
		return nil
	}
	return &_map
}

type Condition struct {
	Attributes map[string]interface{} `json:"attributes"`
}

type ConstantL10n struct {
	Name  map[string]interface{} `json:"name"`
	Error map[string]interface{} `json:"error"`
}

type Constant struct {
	Attributes map[string]interface{} `json:"attributes"`
	Error      string                 `json:"error"`
	L10n       ConstantL10n           `json:"l10n"`
}

//func (s *Schematics) autoTag() {
//	s.Schema.Fields
// add the search tags
//}

func (s *Schematics) Configs() {
	if s.Logging.PrintDebugLogs {
		log.Println("debugger is on")
	}
	if s.Logging.PrintErrorLogs {
		log.Println("error logging is on")
	}
	s.Validators.Logger = s.Logging
	s.Operators.Logger = s.Logging

	if s.Separator == "" {
		s.Separator = "."
	}

}

func (s *Schematics) LoadJsonSchemaFile(path string) error {
	s.Configs()
	content, err := os.ReadFile(path)
	if err != nil {
		s.Logging.ERROR("Failed to load schema file", err)
		return err
	}
	var schema Schema
	err = json.Unmarshal(content, &schema)
	if err != nil {
		s.Logging.ERROR("Failed to unmarshall schema file", err)
		return err
	}
	s.Logging.DEBUG("Schema Loaded From File: ", schema)
	s.Schema = schema
	s.Validators.BasicValidators()
	s.Operators.LoadBasicOperations()
	s.Conditions.BasicConditions()
	if s.Separator == "" {
		s.Separator = "."
	}
	if s.Locale == "" {
		s.Locale = "en"
	}
	return nil
}

func (s *Schematics) LoadMap(schemaMap interface{}) error {
	JSON, err := json.Marshal(schemaMap)
	if err != nil {
		s.Logging.ERROR("Schema should be valid json map[string]interface", err)
		return err
	}
	var schema Schema
	err = json.Unmarshal(JSON, &schema)
	if err != nil {
		s.Logging.ERROR("Invalid Schema", err)
		return err
	}
	s.Logging.DEBUG("Schema Loaded From MAP: ", schema)
	s.Schema = schema
	s.Validators.BasicValidators()
	s.Operators.LoadBasicOperations()
	s.Conditions.BasicConditions()
	if s.Separator == "" {
		s.Separator = "."
	}
	if s.Locale == "" {
		s.Locale = "en"
	}
	return nil
}

// if validators >>> if passed then do *

func fnExists(name string, allValidators map[string]validators.Validator) bool {
	_, exists := allValidators[name]
	if !exists {
		return false
	}
	return true
}

func (f *Field) validateSingleFieldValue(targetID interface{}, value interface{}, allValidators map[string]validators.Validator, db map[string]interface{}, wg *sync.WaitGroup, errChan chan *errorHandler.Error) {
	defer wg.Done()

	var errorMessage errorHandler.Error
	for name, constants := range f.Validators {
		// Early exit if validation name is empty, excluded, or not found in allValidators
		if name == "" || utils.StringInStrings(strings.ToUpper(name), utils.ExcludedValidators) || !fnExists(name, allValidators) {
			continue
		}

		errorMessage.ID = targetID
		errorMessage.Value = value
		errorMessage.Validator = name

		// Set up attributes for validation
		if constants.Attributes == nil {
			constants.Attributes = make(map[string]interface{})
		}
		constants.Attributes["DB"] = db

		// Execute the validator function
		fn := allValidators[name]
		err := fn(value, constants.Attributes)

		// Handle validation errors
		if err != nil {
			// Set custom error message if available
			if constants.Error != "" {
				errorMessage.AddMessage("en", constants.Error)
			}

			// Handle localization (L10n) if present
			if f.L10n != nil {
				for locale, msg := range constants.L10n.Error {
					if msg != nil {
						errorMessage.AddMessage(locale, msg.(string))
					}
				}
				for locale, nameValue := range constants.L10n.Name {
					if nameValue != nil {
						errorMessage.AddL10n(name, locale, nameValue.(string))
					}
				}
			}

			// Send the error message to the error channel
			errChan <- &errorMessage
			return // Exit after handling the first error
		}
	}

	// If no error, return nil to signal success
	errChan <- nil
}

func (f *Field) Validate(allValidators map[string]validators.Validator, id *string, db map[string]interface{}) error {
	if f.Validators == nil {
		return errors.New("no validators defined")
	}
	errorChannel := make(chan *errorHandler.Error)
	var wg sync.WaitGroup
	for targetID, value := range f.Value {
		wg.Add(1)
		go f.validateSingleFieldValue(targetID, value, allValidators, db, &wg, errorChannel)
	}

	go func() {
		wg.Wait()
		close(errorChannel)
	}()

	for vErr := range errorChannel {
		if vErr != nil {
			f.Errors.AddError(f.Target, *vErr)
			f.Status = "failed"
		}
	}
	return nil
}

func (s *Schematics) makeFlat(data map[string]interface{}) *map[string]interface{} {
	if s.Separator == "" {
		s.Separator = "."
	}
	var dMap utils.DataMap
	dMap.FlattenTheMap(data, "", s.Separator)
	s.FlatData = dMap.Data
	return &dMap.Data
}

func (s *Schematics) deflate(data map[string]interface{}) map[string]interface{} {
	unFlatData := utils.DeflateMap(data, s.Separator)
	s.UnFlatData = unFlatData
	return unFlatData
}

func (s *Schematics) AssignData(data map[string]interface{}) error {
	flatData := s.makeFlat(data)
	log.Println("successfully transformed to flat data:", *flatData)
	var fields = make(map[TargetKey]Field)

	if s.Separator == "" {
		s.Separator = "."
	}

	for target, field := range s.Schema.Fields {
		var f = field
		matchingKeys := utils.FindMatchingKeys(*flatData, string(target), s.Separator)
		f.Value = matchingKeys
		log.Println("matching keys ???", matchingKeys, "separator: ", s.Separator)
		if len(matchingKeys) > 0 {
			f.Provided = true
		}

		fields[target] = f
	}

	s.Logging.DEBUG("schema fields: ", utils.InterfaceToJsonString(fields))

	s.Schema.Fields = fields

	return nil
}

func (s *Schematics) Validate(jsonData interface{}) *errorHandler.Errors {
	var baseError errorHandler.Error
	var errs errorHandler.Errors
	baseError.Validator = "validate-object"
	if s == nil {
		baseError.AddMessage("en", "schema not loaded")
		errs.AddError("whole-data", baseError)
		return &errs
	}

	dataBytes, err := json.Marshal(jsonData)
	if err != nil {
		baseError.AddMessage("en", "data is not valid json")
		errs.AddError("whole-data", baseError)
		return &errs
	}

	var obj map[string]interface{}
	var arr []map[string]interface{}
	if err := json.Unmarshal(dataBytes, &obj); err == nil {
		return s.ValidateObject(&obj, nil)
	} else if err := json.Unmarshal(dataBytes, &arr); err == nil {
		return s.ValidateArray(arr)
	} else {
		baseError.AddMessage("en", "invalid format provided for the data, can only be map[string]interface or []map[string]interface")
		errs.AddError("whole-data", baseError)
		return &errs
	}
}

func (s *Schematics) GetValidatedFieldTargets() []string {
	var targets []string
	for target, field := range s.Schema.Fields {
		if field.Provided && field.Errors.HasErrors() == false {
			targets = append(targets, string(target))
		}
	}
	return targets
}

func (f *Field) ConditionalPassage(allConditions conditions.Conditions, schema Schema) bool {
	if len(f.Conditions) > 0 {
		f.logging.DEBUG("inside conditions")
		for name, value := range f.Conditions {
			f.logging.DEBUG("performing conditions", name)
			if fn, ok := allConditions.ConditionFns[name]; ok {
				attrs := value.Attributes
				attrs["schema"] = schema

				if !fn(*f.AsMap(), attrs) {
					return false
				}
			}
		}
	}
	return true
}

func (s *Schematics) ValidateObject(jsonData *map[string]interface{}, id *string) *errorHandler.Errors {
	var errorMessages errorHandler.Errors
	var baseError errorHandler.Error
	flatData := *s.makeFlat(*jsonData)
	err := s.AssignData(flatData)
	if err != nil {
		baseError.AddMessage("en", err.Error())
		errorMessages.AddError("flattening", baseError)
		return &errorMessages
	}

	uniqueID := ""

	if id != nil {
		uniqueID = *id
	}
	db := s.Schema.GetDB(flatData)
	targets := s.GetValidatedFieldTargets()

	for target, field := range s.Schema.Fields {
		if !field.ConditionalPassage(s.Conditions, s.Schema) {
			continue
		}

		if len(field.DependsOn) > 0 {
			missingDependencies := utils.FindUniqueElements(field.DependsOn, targets)
			if len(missingDependencies) > 0 {
				baseError.AddMessage("en", fmt.Sprintf("missing dependencies (%s) for %s", target, strings.Join(missingDependencies, ",")))
				continue
			}
		}

		field.Target = string(target)
		field.logging = s.Logging
		baseError.Validator = "is-required"
		if field.IsRequired && !field.Provided {
			baseError.AddMessage("en", "please provide the value for this required field")
		}
		err := field.Validate(s.Validators.ValidationFns, &uniqueID, db)
		if err != nil {
			baseError.Validator = "common"
			baseError.AddMessage("en", err.Error())
		}
		if field.Errors.HasErrors() {
			errorMessages.MergeErrors(&field.Errors)
		}
	}

	if errorMessages.HasErrors() {
		return &errorMessages
	}
	return nil
}

// GetDB Corrected and completed function
func (s *Schema) GetDB(flatData map[string]interface{}) map[string]interface{} {
	db := s.DB
	for target, field := range s.Fields {
		if field.AddToDB {
			matchingKeys := utils.FindMatchingKeys(flatData, string(target), ".")
			if len(matchingKeys) < 2 {
				mappedKey := utils.GetFirstFromMap(matchingKeys)
				if mappedKey != nil {
					db[string(target)] = mappedKey
				}
			} else if len(matchingKeys) > 0 {
				var values []interface{}
				for _, match := range matchingKeys {
					values = append(values, match)
				}
				db[string(target)] = values
			}
		}
	}
	return db
}

func (s *Schematics) ValidateArray(jsonData []map[string]interface{}) *errorHandler.Errors {
	s.Logging.DEBUG("validating the array")
	var errs errorHandler.Errors
	i := 0
	for _, d := range jsonData {
		var errorMessages *errorHandler.Errors
		var dMap utils.DataMap
		dMap.FlattenTheMap(d, "", s.Separator)
		arrayId, exists := dMap.Data[s.ArrayIdKey]
		if !exists {
			arrayId = fmt.Sprintf("row-%d", i)
			exists = true
		}

		id := arrayId.(string)
		errorMessages = s.ValidateObject(&d, &id)
		if errorMessages.HasErrors() {
			s.Logging.ERROR("has errors", errorMessages.GetStrings("en", "%data\n"))
			errs.MergeErrors(errorMessages)
		}
		i = i + 1
	}

	if errs.HasErrors() {
		return &errs
	}
	return nil
}

// operators

func (f *Field) Operate(value interface{}, allOperations map[string]operators.Op) interface{} {
	for operationName, operationConstants := range f.Operators {
		customValidator, exists := allOperations[operationName]
		if !exists {
			f.logging.ERROR("This operation does not exists in basic or custom operators", operationName)
			return nil
		}
		result := customValidator(value, operationConstants.Attributes)
		if result != nil {
			value = result
		}
	}
	return value
}

func (s *Schematics) Operate(data interface{}) (interface{}, *errorHandler.Errors) {
	var errorMessages errorHandler.Errors
	var baseError errorHandler.Error
	baseError.Validator = "operate-on-schema"
	bytes, err := json.Marshal(data)
	if err != nil {
		s.Logging.ERROR("[operate] error converting the data into bytes", err)
		baseError.AddMessage("en", "data is not valid json")
		errorMessages.AddError("whole-data", baseError)
		return nil, &errorMessages
	}

	dataType, item := utils.IsValidJson(bytes)
	if item == nil {
		s.Logging.ERROR("[operate] error occurred when checking if this data is an array or object")
		baseError.AddMessage("en", "can not convert the data into json")
		errorMessages.AddError("whole-data", baseError)
		return nil, &errorMessages
	}

	if dataType == "object" {
		obj := item.(map[string]interface{})
		results := s.OperateOnObject(obj)
		if results != nil {
			return results, nil
		} else {
			baseError.AddMessage("en", "operation on object unsuccessful")
			errorMessages.AddError("whole-data", baseError)
			return nil, &errorMessages
		}
	} else if dataType == "array" {
		arr := item.([]map[string]interface{})
		results := s.OperateOnArray(arr)
		if results != nil && len(*results) > 0 {
			return results, nil
		} else {
			baseError.AddMessage("en", "operation on array unsuccessful")
			errorMessages.AddError("whole-data", baseError)
			return nil, &errorMessages
		}
	}

	return data, nil
}

func (s *Schematics) OperateOnObject(data map[string]interface{}) *map[string]interface{} {
	data = *s.makeFlat(data)
	for target, field := range s.Schema.Fields {
		matchingKeys := utils.FindMatchingKeys(data, string(target), s.Separator)
		for key, value := range matchingKeys {
			data[key] = field.Operate(value, s.Operators.OpFunctions)
		}
	}
	d := s.deflate(data)
	return &d
}

func (s *Schematics) OperateOnArray(data []map[string]interface{}) *[]map[string]interface{} {
	var obj []map[string]interface{}
	for _, d := range data {
		results := s.OperateOnObject(d)
		obj = append(obj, *results)
	}
	if len(obj) > 0 {
		return &obj
	}
	return nil
}

// General

func (s *Schematics) MergeFields(sc2 *Schematics) *Schematics {
	for target, field := range sc2.Schema.Fields {
		if s.Schema.Fields[target].Type == "" {
			s.Schema.Fields[target] = field
		}
	}
	return s
}
