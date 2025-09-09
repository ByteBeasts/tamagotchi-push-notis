package http

import "errors"

// validateArgs validates the arguments for the http client
func validateArgs(url string, bearer string) error {
	if url == "" {
		return errors.New("url is required")
	}
	if bearer == "" {
		return errors.New("bearer is required")
	}
	return nil
}