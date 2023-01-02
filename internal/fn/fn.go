package fn

import (
	"regexp"
	"strings"
)

func Uniq(s []string) []string {
	k := make(map[string]bool)
	l := []string{}
	for _, i := range s {
		if _, v := k[i]; !v {
			k[i] = true
			l = append(l, i)
		}
	}
	return l
}

func VerifyCountry(input string) bool {
	var re = regexp.MustCompile(`^\w{2}$`)
	ret := re.FindString(input)
	return ret != ""
}

func GetDomains(input string) (result []string) {
	var re = regexp.MustCompile(`^([\.\w\*\-\_]+(\,)?){1,}$`)
	if len(re.FindAllString(input, -1)) > 0 {
		domains := strings.Split(input, ",")
		for _, domain := range domains {
			s := strings.TrimSpace(domain)
			if len(s) > 0 {
				result = append(result, strings.ToLower(domain))
			}
		}
	}
	return result
}

func GetRootDomain(input string) string {
	var re = regexp.MustCompile(`([\.\w\-\_]+){1,2}$`)
	file := strings.TrimLeft(re.FindString(input), ".")
	if file == "" {
		return "cert"
	}
	return file
}
