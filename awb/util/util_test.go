package util_test

import (
	"gitlab.altex.ro/ams/go_ddwcs/awb/util"
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	t.Log("Test with valid email will PASS test")

	//given
	email := "email.123_to@validate.org"

	//when
	if util.IsValidEmail(email) {
		return
	}

	//then
	t.Log("Invalid email: " + email)
	t.FailNow()
}

func TestStringToDateTimeFormat(t *testing.T) {
	t.Log("Test string to date with long date will PASS test")
	date := []string{"2018-11-13T18:29:55.847", "2020-04-08T15:23:46.49"}
	expected := []string{"2018-11-13 18:29", "2020-04-08 15:23"}
	converted := make([]string, 2, 2)

	//when
	for i, value := range date {
		converted[i] = util.StringToDateTimeFormat(value)
	}

	//then
	for i := range converted {
		if converted[i] != expected[i] {
			t.Logf("Expected: %s, Received: %s", expected[i], converted[i])
			t.Fail()
		}
	}
}
