package cmd

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
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

func UpdateBoolOption(key string, args string, defaults string) bool {
	env := os.Getenv(key)
	def := fn.IsBoolString(defaults)
	value := def
	if env != "" {
		value = fn.IsBoolString(env)
	}
	if fn.IsBoolString(args) != def {
		value = fn.IsBoolString(args)
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
	domains := fn.GetDomainsByString(value)
	return domains
}

func SantizeDirPath(key string, args string, defaults string) string {
	value := UpdateStringOption(key, args, defaults)
	if value == defaults {
		return defaults
	}
	s := strings.ToLower(value)
	s = strings.Replace(s, "..", "", -1)
	s = strings.Replace(s, "./", "", -1)
	s = regexp.MustCompile(`[^[:alnum:]\~\-\./]`).ReplaceAllString(s, "")
	return strings.TrimLeft(strings.TrimRight(filepath.Clean(path.Clean(s)), "/"), "/")
}
