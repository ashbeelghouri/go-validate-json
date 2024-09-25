package jsonschematics

import (
	"encoding/json"
	v2 "github.com/ashbeelghouri/jsonschematics/data/v2"
	"github.com/ashbeelghouri/jsonschematics/utils"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"
)

func flatTheMap(data map[string]interface{}) *map[string]interface{} {
	var dataMap utils.DataMap
	dataMap.FlattenTheMap(data, "", ".")
	return &dataMap.Data
}

func FindMatchingKeys(data map[string]interface{}, keyPattern string, separator string) map[string]interface{} {
	matchingKeys := make(map[string]interface{})
	nestedKeys := make(map[string]interface{})
	re := regexp.MustCompile(utils.ConvertKeyToRegex(keyPattern))
	log.Println("regex?", re)
	// Collect all matching keys
	for key, value := range data {
		if strings.HasPrefix(key, keyPattern+separator) {
			log.Println("has prefix?")
			nestedKeys[key] = value
		} else if re.MatchString(key) {
			log.Println("is matching string?", key)
			matchingKeys[key] = value
		}
	}

	if len(nestedKeys) > 0 {
		nest := make(map[string]interface{})

		for key, value := range nestedKeys {
			trimmedKey := strings.TrimPrefix(key, keyPattern+separator)
			nest[trimmedKey] = value
		}

		matchingKeys[keyPattern] = nest
	}

	return matchingKeys
}

func test1() {
	data := map[string]interface{}{
		"name": map[string]interface{}{
			"first": "ashbeel",
			"last":  "ghouri",
		},
	}
	flat := flatTheMap(data)
	log.Println("flat map", flat)

	matchingKeys := FindMatchingKeys(*flat, "name.*", ".")
	log.Println("matching keys: ", matchingKeys)
}

func schemaTest() {
	log.Println("we are inside schema test")
	schematics, err := v2.LoadJsonSchemaFile("test-data/schema/direct/v2/example-1.json")
	if err != nil {
		log.Println("schematics load error ==> ", err)
	}
	schematics.Logging.PrintDebugLogs = true
	schematics.Logging.PrintErrorLogs = true
	content, err := os.ReadFile("test-data/data/direct/example.json")
	if err != nil {
		log.Println(err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Println(err)
	}

	//err = schematics.AssignData(data)
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//log.Println("assigned data", utils.InterfaceToJsonString(schematics.Schema.Fields))

	errs := schematics.Validate(data)
	log.Println("ERRORS: ", errs.GetStrings("en", "%message\n"))
}

func TestV2Validate(t *testing.T) {
	//test1()
	schemaTest()
}
