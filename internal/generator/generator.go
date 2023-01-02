package generator

import (
	"os"
	"strconv"
	"strings"

	"github.com/soulteary/certs-maker/internal/define"
	"github.com/soulteary/certs-maker/internal/fn"
)

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
		define.CERT_DOMAINS = append(define.CERT_DOMAINS, "*")
		define.CERT_DOMAINS = append(define.CERT_DOMAINS, "localhost")
		define.CERT_DOMAINS = fn.Uniq(define.CERT_DOMAINS)
	}

	domains := []string{"[alt_names]"}
	for idx, domain := range define.CERT_DOMAINS {
		id := strconv.Itoa(idx + 1)
		domains = append(domains, "DNS."+id+" = "+domain)
	}
	certDomains := strings.Join(domains, "\n")

	fileName := fn.GetRootDomain(define.CERT_DOMAINS[0])
	if !define.APP_FOR_K8S {
		os.WriteFile("./ssl/"+fileName+".conf", []byte(define.CERT_BASE_INFO+"\n"+certInfo+"\n"+define.CERT_EXTENSIONS+"\n"+certDomains), 0644)
	} else {
		fileName = fileName + ".k8s"
		os.WriteFile("./ssl/"+fileName+".conf", []byte(define.CERT_BASE_INFO+"\n"+certInfo+"\n"+define.CERT_EXTENSIONS_K8S+"\n"+certDomains), 0644)
	}

	scriptTpl := "openssl req -x509 -newkey rsa:2048 -keyout ./ssl/${fileName}.key -out ./ssl/${fileName}.crt -days 3650 -nodes -config ./ssl/${fileName}.conf"
	return strings.ReplaceAll(scriptTpl, "${fileName}", fileName)
}

func TryAdjustPermissions() {
	if define.APP_USER == "" ||
		define.APP_UID == "" ||
		define.APP_GID == "" {
		return
	}
	fn.Execute(`addgroup -g ` + define.APP_GID + ` ` + define.APP_USER)
	fn.Execute(`adduser -g "" -G ` + define.APP_USER + ` -H -D -u ` + define.APP_UID + ` ` + define.APP_USER)
	fn.Execute(`chown -R ` + define.APP_USER + `:` + define.APP_USER + ` ./ssl`)
	fn.Execute(`chmod -R a+r ./ssl`)
}
