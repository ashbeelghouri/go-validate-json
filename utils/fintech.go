package utils

import (
	"fmt"
	"math/big"
	"strings"
	"unicode"
)

var ibanLengths = map[string]int{
	"AL": 28, "AD": 24, "AT": 20, "AZ": 28, "BH": 22, "BE": 16, "BA": 20,
	"BR": 29, "BG": 22, "CR": 21, "HR": 21, "CY": 28, "CZ": 24, "DK": 18,
	"DO": 28, "EE": 20, "FO": 18, "FI": 18, "FR": 27, "GE": 22, "DE": 22,
	"GI": 23, "GR": 27, "GL": 18, "GT": 28, "HU": 28, "IS": 26, "IE": 22,
	"IL": 23, "IT": 27, "JO": 30, "KZ": 20, "XK": 20, "KW": 30, "LV": 21,
	"LB": 28, "LI": 21, "LT": 20, "LU": 20, "MT": 31, "MR": 27, "MU": 30,
	"MC": 27, "MD": 24, "ME": 22, "NL": 18, "MK": 19, "NO": 15, "PK": 24,
	"PS": 29, "PL": 28, "PT": 25, "QA": 29, "RO": 24, "SM": 27, "SA": 24,
	"RS": 22, "SK": 24, "SI": 19, "ES": 24, "SE": 24, "CH": 21, "TN": 24,
	"TR": 26, "UA": 29, "AE": 23, "GB": 22, "VG": 24,
}

func IsValidIBAN(iban string) bool {
	iban = strings.ToUpper(iban)
	if len(iban) < 2 {
		return false
	}
	countryCode := iban[:2]
	if expectedLength, ok := ibanLengths[countryCode]; !ok || len(iban) != expectedLength {
		return false
	}

	// Rearrange the IBAN and convert letters to numbers
	rearranged := iban[4:] + iban[:4]
	var numericIBAN strings.Builder
	for _, ch := range rearranged {
		if unicode.IsLetter(ch) {
			numericIBAN.WriteString(fmt.Sprintf("%d", ch-'A'+10))
		} else {
			numericIBAN.WriteRune(ch)
		}
	}

	// Convert the numeric string to a big integer and perform modulo 97 operation
	numericIBANStr := numericIBAN.String()
	bigIntIBAN, _ := new(big.Int).SetString(numericIBANStr, 10)
	remainder := new(big.Int).Mod(bigIntIBAN, big.NewInt(97))

	return remainder.Cmp(big.NewInt(1)) == 0
}
