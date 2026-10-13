/**
 * Copyright (c) 2021-2025 Su Yang (soulteary)
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

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

func MakeCerts() error {
	expireDays := strings.TrimSpace(define.CERT_EXPIRE_DAYS)
	if expireDays == "" {
		expireDays = define.DEFAULT_EXPIRE_DAYS
	}
	days, err := strconv.ParseUint(expireDays, 10, 31)
	if err != nil || days == 0 {
		return fmt.Errorf("invalid EXPIRE_DAYS %q: expected an integer between 1 and 2147483647", expireDays)
	}
	expireDays = strconv.FormatUint(days, 10)

	baseInfo := GetCertBaseInfo()
	domainList := GetCertDomainList(define.APP_FOR_K8S)

	fileName := GetCertFileNameByDomain(define.CERT_DOMAINS[0], define.APP_FOR_K8S)
	if define.CUSTOM_FILE_NAME != "" {
		fileName = define.CUSTOM_FILE_NAME
	}
	filePath := filepath.Join(define.APP_OUTPUT_DIR, fileName+".conf")

	content := GetCertConfig(baseInfo, domainList, define.APP_FOR_K8S)
	if err := os.WriteFile(filePath, content, define.DEFAULT_MODE); err != nil {
		return fmt.Errorf("write certificate configuration: %w", err)
	}

	if define.APP_FOR_FIREFOX && !define.APP_FOR_K8S {
		if _, err := fn.Execute(GetGeneralExecuteCmds(define.GENERATE_FOR_FF_STEP1, fileName, expireDays)); err != nil {
			return fmt.Errorf("generate Firefox root key: %w", err)
		}
		if _, err := fn.Execute(GetGeneralExecuteCmds(define.GENERATE_FOR_FF_STEP2, fileName, expireDays)); err != nil {
			return fmt.Errorf("generate Firefox root certificate: %w", err)
		}
		if _, err := fn.Execute(GetGeneralExecuteCmds(define.GENERATE_FOR_FF_STEP3, fileName, expireDays)); err != nil {
			return fmt.Errorf("generate Firefox certificate key: %w", err)
		}
		if _, err := fn.Execute(GetGeneralExecuteCmds(define.GENERATE_FOR_FF_STEP4, fileName, expireDays)); err != nil {
			return fmt.Errorf("generate Firefox certificate request: %w", err)
		}
		if _, err := fn.Execute(GetGeneralExecuteCmds(define.GENERATE_FOR_FF_STEP5, fileName, expireDays)); err != nil {
			return fmt.Errorf("sign Firefox certificate: %w", err)
		}
	} else {
		if _, err := fn.Execute(GetGeneralExecuteCmds(define.GENERATE_CMD_TPL, fileName, expireDays)); err != nil {
			return fmt.Errorf("generate certificate: %w", err)
		}
	}

	if _, err := fn.Execute(GetGeneralExecuteCmds(define.CONVERT_CRT_TO_DER, fileName, expireDays)); err != nil {
		return fmt.Errorf("convert certificate to DER: %w", err)
	}
	if _, err := fn.Execute(GetGeneralExecuteCmds(define.CONVERT_KEY_TO_DER, fileName, expireDays)); err != nil {
		return fmt.Errorf("convert certificate key to DER: %w", err)
	}
	return nil
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
	extensions := define.CERT_EXTENSIONS
	if isK8s {
		extensions = define.CERT_EXTENSIONS_K8S
	} else if define.APP_FOR_FIREFOX {
		extensions = define.CERT_EXTENSIONS_FIREFOX
	}
	return []byte(fmt.Sprintf("%s\n%s\n%s\n%s\n", define.CERT_BASE_INFO, info, extensions, domain))
}

func GetGeneralExecuteCmds(commandTpl string, output string, exipre string) string {
	cmd := strings.ReplaceAll(commandTpl, define.GENERATE_FILE_PLACEHOLDER, fmt.Sprintf("%s/%s", define.APP_OUTPUT_DIR, output))
	return strings.ReplaceAll(cmd, define.GENERATE_EXPIRE_DAYS_PLACEHOLDER, exipre)
}
