package utils

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/url"
	"regexp"
	"strings"
)

var httpsCheck = regexp.MustCompile(`^https?://`)
var httpCheck = regexp.MustCompile(`^http?://`)

// GetKeysFromMapStringInterface parses keys from a map[string]interface
func GetKeysFromMapStringInterface(m map[string]interface{}) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

// GetKeysFromMapString returns the values of a map[string]string keys as a slice of strings
func GetKeysFromMapString(m map[string]string) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

// isHTTPS determines whether the string url has an https scheme prefix
func isHTTPS(url string) bool {
	return len(httpsCheck.FindAllString(url, -1)) != 0
}

// isHTTP determines whether the string url has an http scheme prefix
func isHTTP(url string) bool {
	return len(httpCheck.FindAllString(url, -1)) != 0
}

// ValidateURL accepts a path and validates it as a url
func ValidateURL(path string) (string, error) {
	if !isHTTP(path) && !isHTTPS(path) {
		if strings.Contains(path, "localhost") {
			path = fmt.Sprintf("http://%s", path)
		} else {
			path = fmt.Sprintf("https://%s", path)
		}
	}

	u, err := url.ParseRequestURI(path)
	if err != nil {
		PrintError(err)
		return "", err
	}

	valid := govalidator.IsRequestURL(u.String())
	if !valid {
		return "", fmt.Errorf("%v is not a valid URL", u.String())
	}

	return u.String(), nil
}
