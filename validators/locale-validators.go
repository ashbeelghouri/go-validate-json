package validators

import (
	"encoding/json"
	"errors"
	"github.com/DScale-io/jsonschematics/validators/archives"
	"os"
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

func IsCountryValid(i interface{}, _ map[string]interface{}) error {
	userCountry := i.(string)
	countries := archives.GetCountries()
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
