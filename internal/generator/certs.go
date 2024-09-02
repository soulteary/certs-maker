package generator

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/soulteary/certs-maker/internal/define"
	"github.com/soulteary/certs-maker/internal/fn"
)

func MakeCerts() {
	baseInfo := GetCertBaseInfo()
	domainList := GetCertDomainList(define.APP_FOR_K8S)

	fileName := GetCertFileNameByDomain(define.CERT_DOMAINS[0], define.APP_FOR_K8S)
	filePath := filepath.Join(define.APP_OUTPUT_DIR, fileName+".conf")

	content := GetCertConfig(baseInfo, domainList, define.APP_FOR_K8S)
	os.WriteFile(filePath, content, define.DEFAULT_MODE)

	expireDays := define.DEFAULT_EXPIRE_DAYS

	if define.APP_FOR_FIREFOX && !define.APP_FOR_K8S {
		fn.Execute(GetGeneralExecuteCmds(define.GENERATE_FOR_FF_STEP1, fileName, expireDays))
		fn.Execute(GetGeneralExecuteCmds(define.GENERATE_FOR_FF_STEP2, fileName, expireDays))
		fn.Execute(GetGeneralExecuteCmds(define.GENERATE_FOR_FF_STEP3, fileName, expireDays))
		fn.Execute(GetGeneralExecuteCmds(define.GENERATE_FOR_FF_STEP4, fileName, expireDays))
		fn.Execute(GetGeneralExecuteCmds(define.GENERATE_FOR_FF_STEP5, fileName, expireDays))
	} else {
		fn.Execute(GetGeneralExecuteCmds(define.GENERATE_CMD_TPL, fileName, expireDays))
	}

	fn.Execute(GetGeneralExecuteCmds(define.CONVERT_CRT_TO_DER, fileName, expireDays))
	fn.Execute(GetGeneralExecuteCmds(define.CONVERT_KEY_TO_DER, fileName, expireDays))
}

func GetCertBaseInfo() string {
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

func GetCertDomainList(isK8s bool) string {
	if isK8s {
		define.CERT_DOMAINS = append(define.CERT_DOMAINS, "*", "localhost")
		define.CERT_DOMAINS = fn.GetUniqDomains(define.CERT_DOMAINS)
	}

	domains := []string{"[alt_names]"}
	for idx, domain := range define.CERT_DOMAINS {
		id := strconv.Itoa(idx + 1)
		if net.ParseIP(domain) == nil {
			domains = append(domains, "DNS."+id+" = "+domain)
		} else {
			domains = append(domains, "IP."+id+" = "+domain)
		}
	}
	return strings.Join(domains, "\n")
}

func GetCertFileNameByDomain(domain string, isK8s bool) string {
	s := fn.GetDomainName(domain)
	if isK8s {
		s += ".k8s"
	}
	return s
}

func GetCertConfig(info string, domain string, isK8s bool) []byte {
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

func GetGeneralExecuteCmds(commandTpl string, output string, exipre string) string {
	cmd := strings.ReplaceAll(commandTpl, define.GENERATE_FILE_PLACEHOLDER, fmt.Sprintf("%s/%s", define.APP_OUTPUT_DIR, output))
	return strings.ReplaceAll(cmd, define.GENERATE_EXPIRE_DAYS_PLACEHOLDER, exipre)
}
