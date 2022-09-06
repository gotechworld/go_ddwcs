package util

import (
	"regexp"
	"strings"
	"time"
)

/**
 * Email validation method
 */
func IsValidEmail(email string) bool {
	rule := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if rule.MatchString(email) {
		return true
	}

	return false
}

/**
 * Create a string key from order params
 * customer_id or customer_email will contains an empty string
 */
func CreateKeyFromOrderParams(params map[string]string) string {
	key := params["increment_id"] + ":"
	key += params["customer_id"]
	key += params["customer_email"]

	return key
}

/**
 * Create a string in (Y-m-d H:i) format from a RFC3999
 * It returns the current date if the string does not match with the regex
 */
func StringToDateTimeFormat(datetime string) string {
	rule := regexp.MustCompile(`(?i)([1][9][0-9]{2}|[2][0-9]{3})-([0][1-9]|[1][0-2])-([0][1-9]|[1][0-9]|[2][0-9]|[3][0-1]).([0-1][0-9]|[2][0-3]):([0-5][0-9]):([0-5][0-9])(.*)`)
	if rule.MatchString(datetime) {
		datetime = strings.Replace(datetime, "T", " ", 1)
		return datetime[0:16]
	}

	return time.Now().Format("2006-01-02 15:04")
}
