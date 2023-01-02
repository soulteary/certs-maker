package fn

import (
	"reflect"
	"strings"
)

func IsBoolString(input string) bool {
	s := strings.TrimSpace(strings.ToLower(input))
	if s == "true" || s == "1" || s == "on" {
		return true
	}
	return false
}

func IsNotEmptyAndNotDefaultString(value string, defaults string) bool {
	return value != "" && value != defaults
}

func IsDomainListStringMatch(value []string, defaults string) bool {
	return reflect.DeepEqual(value, GetDomains(defaults))
}
