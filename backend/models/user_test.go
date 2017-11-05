package models

import (
	"testing"
)

//
// Validate an email address
//
func TestValidateEmailAddress(t *testing.T) {

	testData := map[string]bool{
		"spicer@options.cafe": true,
		"spicer matthews":     false,
		"@woot.com":           false,
		"me@example.com":      true,
	}

	for email, isReal := range testData {

		result := ValidateEmailAddress(email)

		if isReal && (result != nil) {
			t.Errorf("%s did not pass.", email)
		} else if (!isReal) && (result == nil) {
			t.Errorf("%s did not pass.", email)
		}

	}
}

//
// Test GenerateRandomBytes returns
//
func TestGenerateRandomBytes(t *testing.T) {

	randString, _ := GenerateRandomBytes(10)

	if len(randString) != 10 {
		t.Errorf("The random bytes was %d chars instead of 10.", len(randString))
	}

}

//
// Test GenerateRandomString returns
//
func TestGenerateRandomString(t *testing.T) {

	randString, _ := GenerateRandomString(10)

	if len(randString) != 10 {
		t.Errorf("The random string of %s was %d chars instead of 10.", randString, len(randString))
	}

}

/* End File */
