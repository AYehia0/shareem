package validation

import (
	"net/url"
)

// responsible for request validation

// checks if the input is a valid URL
func IsValidURL(input string) bool {
	u, err := url.Parse(input)

	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return true
}
