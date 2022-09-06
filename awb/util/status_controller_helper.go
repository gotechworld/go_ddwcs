package util

import (
	"errors"
	"net/url"
	"strings"
)

type StatusControllerHelper struct{}

/**
 * ExtractCustomerFromQueryString
 * @param queryString url.Values
 * @return string
 * @desc - treat the special case about customer - we accept either customer_id or customer_email
 */
func (cpv *StatusControllerHelper) ExtractCustomerFromQueryString(queryString *url.Values) string {
	customer := queryString.Get("customer_id")
	if customer == "" {
		customer = queryString.Get("customer_email")
	}

	return customer
}

/**
 * ValidateParams
 * @param queryString *url.Values
 * @param whitelist []string
 * @return error|nil
 */
func (cpv *StatusControllerHelper) ValidateParams(queryString *url.Values, whitelist []string) error {
	var isValid bool = false
	for _, elem := range whitelist {

		if strings.Index(elem, "||") != -1 {
			isValid = false
			for _, elem := range strings.Split(elem, "||") {
				if queryString.Get(elem) != "" {
					isValid = true
					break
				}
			}
		} else if queryString.Get(elem) != "" {
			isValid = true
		}

		if !isValid {
			return errors.New("invalid parameter: " + elem)
		}
	}

	return nil
}
