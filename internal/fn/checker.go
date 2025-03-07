/**
 * Copyright (c) 2021-2025 Su Yang (soulteary)
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package fn

import (
	"reflect"
	"regexp"
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
	return reflect.DeepEqual(value, GetDomainsByString(defaults))
}

func IsVaildCountry(input string) bool {
	return regexp.MustCompile(`^\w{2}$`).FindString(input) != ""
}

func IsStrInArray(arr []string, s string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}
