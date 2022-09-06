package util_test

import (
	"gitlab.altex.ro/ams/go_ddwcs/awb/util"
	"net/url"
	"testing"
)

func TestStatusControllerHelper_ValidateParams(t *testing.T) {
	t.Log("Check for parameters with all values available - simple case will pass")

	//given
	statusWhitelist := []string{"increment_id", "customer_id"}
	queryString := url.Values{
		"increment_id": []string{"test_increment_id"},
		"customer_id":  []string{"123123"},
	}
	helper := &util.StatusControllerHelper{}

	//when
	err := helper.ValidateParams(&queryString, statusWhitelist)

	//then
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
}

func TestStatusControllerHelper_ValidateParamsOrCase(t *testing.T) {
	t.Log("Check for parameters with all values available - OR case will pass")

	//given
	statusWhitelist := []string{"increment_id", "customer_id||customer_email"}
	queryString := url.Values{
		"increment_id": []string{"test_increment_id"},
		"customer_id":  []string{"123123"},
	}
	helper := &util.StatusControllerHelper{}

	//when
	err := helper.ValidateParams(&queryString, statusWhitelist)

	//then
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
}

func TestStatusControllerHelper_ValidateParamsOrCaseMissingValues(t *testing.T) {
	t.Log("Check for parameters with all values available - OR case - query string missing params will return error message")

	//given
	statusWhitelist := []string{"increment_id", "customer_id||customer_email"}
	queryString := url.Values{
		"increment_id": []string{"test_increment_id"},
	}
	helper := &util.StatusControllerHelper{}

	//when
	err := helper.ValidateParams(&queryString, statusWhitelist)

	//then
	if err != nil {
		t.Log(err.Error())
		return
	}

	t.FailNow()
}
