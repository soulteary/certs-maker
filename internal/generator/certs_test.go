package generator_test

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/soulteary/certs-maker/internal/define"
	"github.com/soulteary/certs-maker/internal/generator"
)

func TestGetCertBaseInfo(t *testing.T) {

	define.CERT_COUNTRY = "ABC"
	define.CERT_STATE = "BCD"
	define.CERT_LOCALITY = "CDE"
	define.CERT_ORGANIZATION = "DEF"
	define.CERT_ORGANIZATIONAL_UNIT = "EFG"
	define.CERT_COMMON_NAME = "FGH"

	info := generator.GetCertBaseInfo()

	check1 := strings.Contains(info, fmt.Sprintf("C  = %s", define.CERT_COUNTRY))
	if !check1 {
		t.Fatal("Test GetCertBaseInfo failed")
	}

	check2 := strings.Contains(info, fmt.Sprintf("ST = %s", define.CERT_STATE))
	if !check2 {
		t.Fatal("Test GetCertBaseInfo failed")
	}

	check3 := strings.Contains(info, fmt.Sprintf("L  = %s", define.CERT_LOCALITY))
	if !check3 {
		t.Fatal("Test GetCertBaseInfo failed")
	}

	check4 := strings.Contains(info, fmt.Sprintf("O  = %s", define.CERT_ORGANIZATION))
	if !check4 {
		t.Fatal("Test GetCertBaseInfo failed")
	}

	check5 := strings.Contains(info, fmt.Sprintf("OU = %s", define.CERT_ORGANIZATIONAL_UNIT))
	if !check5 {
		t.Fatal("Test GetCertBaseInfo failed")
	}

	check6 := strings.Contains(info, fmt.Sprintf("CN = %s", define.CERT_COMMON_NAME))
	if !check6 {
		t.Fatal("Test GetCertBaseInfo failed")
	}
}

func TestGetCertDomainList(t *testing.T) {
	ret := generator.GetCertDomainList(false)
	if ret != "[alt_names]" {
		t.Fatal("Test GetCertDomainList failed")
	}

	define.CERT_DOMAINS = append(define.CERT_DOMAINS, "abc.com")
	ret = generator.GetCertDomainList(false)
	if !strings.Contains(ret, "DNS.1 = abc.com") {
		t.Fatal("Test GetCertDomainList failed")
	}

	define.CERT_DOMAINS = append(define.CERT_DOMAINS, "127.0.0.1")
	ret = generator.GetCertDomainList(false)
	if !strings.Contains(ret, "IP.2 = 127.0.0.1") {
		t.Fatal("Test GetCertDomainList failed")
	}

	ret = generator.GetCertDomainList(true)
	if !strings.Contains(ret, "*") || !strings.Contains(ret, "localhost") {
		t.Fatal("Test GetCertDomainList failed")
	}
}

func TestGetCertFileNameByDomain(t *testing.T) {
	ret := generator.GetCertFileNameByDomain("abc.com", false)
	if ret != "abc.com" {
		t.Fatal("test GetCertFileNameByDomain failed")
	}
	ret = generator.GetCertFileNameByDomain("abc.com", true)
	if ret != "abc.com.k8s" {
		t.Fatal("test GetCertFileNameByDomain failed")
	}
}

func TestGetCertConfig(t *testing.T) {
	ret := string(generator.GetCertConfig("this-is-info", "here-is-domain", false))
	if !strings.Contains(ret, "this-is-info") {
		t.Fatal("test GetCertConfig failed")
	}
	if !strings.Contains(ret, "here-is-domain") {
		t.Fatal("test GetCertConfig failed")
	}
	if strings.Contains(ret, "basicConstraints = CA:FALSE") {
		t.Fatal("test GetCertConfig failed")
	}

	ret = string(generator.GetCertConfig("this-is-info", "here-is-domain", true))
	if !strings.Contains(ret, "this-is-info") {
		t.Fatal("test GetCertConfig failed")
	}
	if !strings.Contains(ret, "here-is-domain") {
		t.Fatal("test GetCertConfig failed")
	}
	if !strings.Contains(ret, "basicConstraints = CA:FALSE") {
		t.Fatal("test GetCertConfig failed")
	}
}

func TestGetGeneralExecuteCmds(t *testing.T) {
	ret := generator.GetGeneralExecuteCmds(define.GENERATE_CMD_TPL, "abc")
	if ret != "openssl req -x509 -newkey rsa:2048 -keyout /abc.key -out /abc.crt -days 3650 -nodes -config /abc.conf" {
		t.Fatal("test GetGeneralExecuteCmds failed")
	}
}

func TestMakeCerts(t *testing.T) {
	// test for common use
	define.APP_FOR_K8S = false
	define.CERT_DOMAINS = append(define.CERT_DOMAINS, "abc.com")
	define.APP_OUTPUT_DIR = "./"

	generator.MakeCerts()

	fileGenerated := []string{
		"./abc.com.conf",
		"./abc.com.key",
	}

	for _, file := range fileGenerated {
		if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
			t.Fatal("test MakeCerts failed")
		}
		os.Remove(file)
	}

	// test for k8s use
	define.APP_FOR_K8S = true
	define.CERT_DOMAINS = append(define.CERT_DOMAINS, "abc.com")
	define.APP_OUTPUT_DIR = "./"

	generator.MakeCerts()

	fileGenerated = []string{
		"./abc.com.k8s.conf",
		"./abc.com.k8s.key",
	}
	for _, file := range fileGenerated {
		if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
			t.Fatal("test MakeCerts failed")
		}
		os.Remove(file)
	}

	// test for firefox use
	define.APP_FOR_K8S = false
	define.APP_FOR_FIREFOX = true
	define.CERT_DOMAINS = append(define.CERT_DOMAINS, "abc.com")
	define.APP_OUTPUT_DIR = "./"

	generator.MakeCerts()

	fileGenerated = []string{
		"./abc.com.conf",
		"./abc.com.key",
		"./abc.com.rootCA.key",
	}
	for _, file := range fileGenerated {
		if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
			t.Fatal("test MakeCerts failed")
		}
		os.Remove(file)
	}
}
