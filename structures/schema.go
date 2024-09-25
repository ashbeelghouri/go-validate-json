package structures

import (
	"github.com/ashbeelghouri/jsonschematics/errorHandler"
	"github.com/ashbeelghouri/jsonschematics/utils"
)

type TargetKey string

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

type Condition struct {
	Action     string                 `json:"action"`
	attributes map[string]interface{} `json:"attributes"`
}

type ConditionalAction struct {
	Success []string `json:"success"`
	Error   []string `json:"error"`
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
