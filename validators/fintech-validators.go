package validators

import (
	"errors"
	"github.com/DScale-io/jsonschematics/utils"
)

func IsValidIBAN(i interface{}, _ map[string]interface{}) error {
	iban := i.(string)
	if !utils.IsValidIBAN(iban) {
		return errors.New("invalid IBAN provided")
	}
	return nil
}
