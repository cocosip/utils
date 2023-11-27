package util

import "net/url"

func IsURL(rowURL string) bool {
	_, err := url.Parse(rowURL)
	return err == nil
}
