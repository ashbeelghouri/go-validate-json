package validators

import (
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strings"
)

func getCountriesList() *map[string]string {
	var list map[string]string
	content, err := os.ReadFile("archives/countriesList.json")
	if err != nil {
		return nil
	}
	err = json.Unmarshal(content, &list)
	if err != nil {
		return nil
	}
	return &list
}

func IsValidEmiratesID(id string) bool {
	// Regular expression to match the Emirates ID format
	re := regexp.MustCompile(`^784-\d{4}-\d{7}-\d$`)
	return re.MatchString(id)
}

func IsCountryValid(i interface{}, _ map[string]interface{}) error {
	userCountry := i.(string)
	countries := getCountriesList()

	for code, country := range *countries {
		uc := strings.ToLower(userCountry)
		c := strings.ToLower(country)
		cd := strings.ToLower(code)
		if uc != c && uc != cd {
			return errors.New("this is an invalid country")
		}
	}

	return nil
}

func IsEmiratesIDValid(i interface{}, _ map[string]interface{}) error {
	emiratesID := i.(string)
	if !IsValidEmiratesID(emiratesID) {
		return errors.New("invalid emirates id provided")
	}
	return nil
}
