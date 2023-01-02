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

func GenerateConfFile() string {
	certInfo := strings.Join(
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

	if define.APP_FOR_K8S {
		define.CERT_DOMAINS = append(define.CERT_DOMAINS, "*", "localhost")
		define.CERT_DOMAINS = fn.GetUniqDomains(define.CERT_DOMAINS)
	}

	domains := []string{"[alt_names]"}
	for idx, domain := range define.CERT_DOMAINS {
		id := strconv.Itoa(idx + 1)
		domains = append(domains, "DNS."+id+" = "+domain)
	}
	certDomains := strings.Join(domains, "\n")

	fileName := fn.GetDomainName(define.CERT_DOMAINS[0])
	if !define.APP_FOR_K8S {
		confPath := filepath.Join(define.APP_DIR, fileName+".conf")
		os.WriteFile(confPath, []byte(define.CERT_BASE_INFO+"\n"+certInfo+"\n"+define.CERT_EXTENSIONS+"\n"+certDomains), 0644)
	} else {
		fileName = fileName + ".k8s"
		confPath := filepath.Join(define.APP_DIR, fileName+".conf")
		os.WriteFile(confPath, []byte(define.CERT_BASE_INFO+"\n"+certInfo+"\n"+define.CERT_EXTENSIONS_K8S+"\n"+certDomains), 0644)
	}

	scriptTpl := "openssl req -x509 -newkey rsa:2048 -keyout ${file}.key -out ${file}.crt -days 3650 -nodes -config ${file}.conf"
	return strings.ReplaceAll(scriptTpl, "${file}", fmt.Sprintf("%s/%s", define.APP_DIR, fileName))
}

func TryAdjustPermissions() {
	if define.APP_USER == "" ||
		define.APP_UID == "" ||
		define.APP_GID == "" {
		return
	}
	fn.Execute(`addgroup -g ` + define.APP_GID + ` ` + define.APP_USER)
	fn.Execute(`adduser -g "" -G ` + define.APP_USER + ` -H -D -u ` + define.APP_UID + ` ` + define.APP_USER)
	fn.Execute(`chown -R ` + define.APP_USER + `:` + define.APP_USER + ` ` + define.APP_DIR)
	fn.Execute(`chmod -R a+r ` + define.APP_DIR)
}
