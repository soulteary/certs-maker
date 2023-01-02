package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/soulteary/certs-maker/internal/define"
	"github.com/soulteary/certs-maker/internal/fn"
)

func Generate() {
	shell := GenerateConfFile()
	fn.Execute(shell)
	TryAdjustPermissions()
}

func GetCertSettings() string {
	return strings.Join(
		[]string{
			"[cert_distinguished_name]",
			"C  = " + define.CERT_COUNTRY,
			"ST = " + define.CERT_STATE,
			"L  = " + define.CERT_LOCALITY,
			"O  = " + define.CERT_ORGANIZATION,
			"OU = " + define.CERT_ORGANIZATIONAL_UNIT,
			"CN = " + define.CERT_COMMON_NAME,
		}, "\n",
	)
}

func GetCertDomains(isK8s bool) string {
	if isK8s {
		define.CERT_DOMAINS = append(define.CERT_DOMAINS, "*", "localhost")
		define.CERT_DOMAINS = fn.GetUniqDomains(define.CERT_DOMAINS)
	}

	domains := []string{"[alt_names]"}
	for idx, domain := range define.CERT_DOMAINS {
		id := strconv.Itoa(idx + 1)
		domains = append(domains, "DNS."+id+" = "+domain)
	}
	return strings.Join(domains, "\n")
}

func GetCertFullConfig(info string, domain string, isK8s bool) []byte {
	if !isK8s {
		return []byte(
			fmt.Sprintf("%s\n%s\n%s\n%s\n", define.CERT_BASE_INFO, info, define.CERT_EXTENSIONS, domain),
		)
	} else {
		return []byte(
			fmt.Sprintf("%s\n%s\n%s\n%s\n", define.CERT_BASE_INFO, info, define.CERT_EXTENSIONS_K8S, domain),
		)
	}
}

func GetCertFileNameByDomain(domain string, isK8s bool) string {
	s := fn.GetDomainName(domain)
	if isK8s {
		s += ".k8s"
	}
	return s
}

func GenerateConfFile() string {
	certInfo := GetCertSettings()
	certDomains := GetCertDomains(define.APP_FOR_K8S)

	fileName := GetCertFileNameByDomain(define.CERT_DOMAINS[0], define.APP_FOR_K8S)
	confPath := filepath.Join(define.APP_OUTPUT_DIR, fileName+".conf")

	content := GetCertFullConfig(certInfo, certDomains, define.APP_FOR_K8S)
	os.WriteFile(confPath, content, define.DEFAULT_MODE)

	scriptTpl := "openssl req -x509 -newkey rsa:2048 -keyout ${file}.key -out ${file}.crt -days 3650 -nodes -config ${file}.conf"
	return strings.ReplaceAll(scriptTpl, "${file}", fmt.Sprintf("%s/%s", define.APP_OUTPUT_DIR, fileName))
}

func TryAdjustPermissions() {
	if define.APP_USER == "" ||
		define.APP_UID == "" ||
		define.APP_GID == "" {
		return
	}

	cmdCreateGroup := fmt.Sprintf(`addgroup -g %s %s`, define.APP_GID, define.APP_USER)
	cmdCreateUser := fmt.Sprintf(`adduser -g "" -G %s -H -D -u %s %s`, define.APP_USER, define.APP_UID, define.APP_USER)
	cmdChangeOwner := fmt.Sprintf(`chown -R %s:%s %s`, define.APP_USER, define.APP_USER, define.APP_OUTPUT_DIR)
	cmdChangeMod := fmt.Sprintf(`chmod -R a+r %s`, define.APP_OUTPUT_DIR)

	fn.Execute(fmt.Sprintf("%s\n%s\n%s\n%s\n", cmdCreateGroup, cmdCreateUser, cmdChangeOwner, cmdChangeMod))
}
