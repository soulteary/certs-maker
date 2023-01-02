package cmd

import (
	"os"
	"strings"

	"github.com/soulteary/certs-maker/internal/fn"
)

func UpdateStringOption(key string, args string, defaults string) string {
	env := os.Getenv(key)
	str := defaults
	if fn.IsNotEmptyAndNotDefaultString(env, defaults) {
		str = env
	}
	if fn.IsNotEmptyAndNotDefaultString(args, defaults) {
		str = args
	}
	return strings.TrimSpace(str)
}

func UpdateBoolOption(key string, args bool, defaults bool) bool {
	env := os.Getenv(key)
	value := defaults
	if env != "" {
		value = fn.IsBoolString(env)
	}
	if args != defaults {
		value = args
	}
	return value
}

func UpdateCountryOption(key string, args string, defaults string) string {
	value := UpdateStringOption(key, args, defaults)
	if fn.IsVaildCountry(value) {
		return strings.ToUpper(value)
	}
	return defaults
}

func UpdateDomainOption(key string, args string, defaults string) []string {
	value := UpdateStringOption(key, args, defaults)
	domains := fn.GetDomains(value)
	return domains
}
