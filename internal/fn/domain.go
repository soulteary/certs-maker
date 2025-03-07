/**
 * Copyright (c) 2021-2025 Su Yang (soulteary)
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package fn

import (
	"fmt"
	"net"
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
	domainPattern := `(\*\.)?([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}`
	ipv4Pattern := `((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`
	ipv6Pattern := `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	combinedPattern := fmt.Sprintf(`^(%s|%s|%s)$`, domainPattern, ipv4Pattern, ipv6Pattern)
	re := regexp.MustCompile(combinedPattern)

	addresses := strings.Split(input, ",")

	for _, addr := range addresses {
		addr = strings.TrimSpace(addr)
		if len(addr) == 0 {
			continue
		}

		if re.MatchString(addr) {
			if strings.Contains(addr, ":") || net.ParseIP(addr) != nil {
				result = append(result, addr)
			} else {
				result = append(result, strings.ToLower(addr))
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
