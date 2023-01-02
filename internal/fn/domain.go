package fn

import (
	"regexp"
	"strings"
)

func GetUniqDomains(s []string) []string {
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

func GetDomainsByString(input string) (result []string) {
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
	return GetUniqDomains(result)
}

func GetDomainName(input string) string {
	var re = regexp.MustCompile(`([\.\w\-\_]+){1,2}$`)
	domain := strings.TrimLeft(re.FindString(input), ".")
	return domain
}
