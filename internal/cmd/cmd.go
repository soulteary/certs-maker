package cmd

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"cert-maker/internal/define"
	"cert-maker/internal/fn"
)

func parseEnvInputs() (cert define.CERT) {
	country := os.Getenv(define.ENV_KEY_COUNTRY)
	state := os.Getenv(define.ENV_KEY_STATE)
	locality := os.Getenv(define.ENV_KEY_LOCALITY)
	organization := os.Getenv(define.ENV_KEY_ORGANIZATION)
	organizationalUnit := os.Getenv(define.ENV_KEY_ORGANIZATION_UNIT)
	commonName := os.Getenv(define.ENV_KEY_COMMON_NAME)
	domains := os.Getenv(define.ENV_KEY_DOMAINS)
	forK8S := os.Getenv(define.ENV_KEY_FOR_K8S)
	user := os.Getenv(define.ENV_KEY_USER)
	uid := os.Getenv(define.ENV_KEY_UID)
	gid := os.Getenv(define.ENV_KEY_GID)

	return createCertConfig(country, state, locality, organization, organizationalUnit, commonName, domains, forK8S, user, uid, gid)
}

func parseCliInputs() (cert define.CERT) {
	var country string
	flag.StringVar(&country, define.ENV_KEY_COUNTRY, define.DEFAULT_COUNTRY, "Country Name")

	var state string
	flag.StringVar(&state, define.ENV_KEY_STATE, define.DEFAULT_STATE, "State Or Province Name")

	var locality string
	flag.StringVar(&locality, define.ENV_KEY_LOCALITY, define.DEFAULT_LOCALITY, "Locality Name")

	var organization string
	flag.StringVar(&organization, define.ENV_KEY_ORGANIZATION, define.DEFAULT_ORGANIZATION, "Organization Name")

	var organizationalUnit string
	flag.StringVar(&organizationalUnit, define.ENV_KEY_ORGANIZATION_UNIT, define.DEFAULT_ORGANIZATIONAL_UNIT, "Organizational Unit Name")

	var commonName string
	flag.StringVar(&commonName, define.ENV_KEY_COMMON_NAME, define.DEFAULT_COMMON_NAME, "Common Name")

	var domains string
	flag.StringVar(&domains, define.ENV_KEY_DOMAINS, define.DEFAULT_DOMAINS, "Domains")

	var forK8S string
	flag.StringVar(&forK8S, define.ENV_KEY_FOR_K8S, define.DEFAULT_FORK8S, "FOR K8S")

	var user string
	flag.StringVar(&user, define.ENV_KEY_USER, "", "File Owner Username")

	var uid string
	flag.StringVar(&uid, define.ENV_KEY_UID, "", "File Owner UID")

	var gid string
	flag.StringVar(&gid, define.ENV_KEY_GID, "", "File Owner GID")

	flag.Parse()

	return createCertConfig(country, state, locality, organization, organizationalUnit, commonName, domains, forK8S, user, uid, gid)
}

func MergeUserInputs() define.CERT {
	base := parseEnvInputs()
	cli := parseCliInputs()

	if cli.Country != define.DEFAULT_COUNTRY {
		base.Country = cli.Country
	}
	if cli.State != define.DEFAULT_STATE {
		base.State = cli.State
	}
	if cli.Locality != define.DEFAULT_LOCALITY {
		base.Locality = cli.Locality
	}
	if cli.Organization != define.DEFAULT_ORGANIZATION {
		base.Organization = cli.Organization
	}
	if cli.OrganizationalUnit != define.DEFAULT_ORGANIZATIONAL_UNIT {
		base.OrganizationalUnit = cli.OrganizationalUnit
	}
	if cli.CommonName != define.DEFAULT_COMMON_NAME {
		base.CommonName = cli.CommonName
	}
	if !reflect.DeepEqual(cli.Domains, fn.GetDomains(define.DEFAULT_DOMAINS)) {
		base.Domains = cli.Domains
	}
	if cli.ForK8S != base.ForK8S {
		base.ForK8S = cli.ForK8S
	}
	if cli.OwnUser != "" && cli.OwnUID != "" && cli.OwnGID != "" {
		if cli.OwnUser != base.OwnUser && cli.OwnUID != base.OwnUID && cli.OwnGID != base.OwnGID {
			base.OwnUser = cli.OwnUser
			base.OwnUID = cli.OwnUID
			base.OwnGID = cli.OwnGID
		}
	}

	return base
}

func createCertConfig(country string, state string, locality string, organization string, organizationalUnit string, commonName string, domains string, forK8S string, user string, uid string, gid string) (cert define.CERT) {
	country = strings.TrimSpace(country)
	if len(country) > 0 {
		if fn.VerifyCountry(country) {
			cert.Country = strings.ToUpper(country)
		} else {
			fmt.Println("wrong country name, set to default value:", define.DEFAULT_COUNTRY)
		}
	} else {
		cert.Country = define.DEFAULT_COUNTRY
	}

	state = strings.TrimSpace(state)
	if len(state) > 0 {
		cert.State = strings.ToUpper(state)
	} else {
		cert.State = define.DEFAULT_STATE
	}

	locality = strings.TrimSpace(locality)
	if len(locality) > 0 {
		cert.Locality = strings.ToUpper(locality)
	} else {
		cert.Locality = define.DEFAULT_LOCALITY
	}

	organization = strings.TrimSpace(organization)
	if len(organization) > 0 {
		cert.Organization = organization
	} else {
		cert.Organization = define.DEFAULT_ORGANIZATION
	}

	organizationalUnit = strings.TrimSpace(organizationalUnit)
	if len(organization) > 0 {
		cert.OrganizationalUnit = organizationalUnit
	} else {
		cert.OrganizationalUnit = define.DEFAULT_ORGANIZATIONAL_UNIT
	}

	commonName = strings.TrimSpace(commonName)
	if len(commonName) > 0 {
		cert.CommonName = commonName
	} else {
		cert.CommonName = define.DEFAULT_COMMON_NAME
	}

	domainsInput := strings.TrimSpace(domains)
	if len(domainsInput) > 0 {
		userDomains := fn.GetDomains(domainsInput)
		if len(userDomains) == 0 {
			userDomains = fn.GetDomains(define.DEFAULT_DOMAINS)
			fmt.Println("wrong domains, set to default value:", define.DEFAULT_DOMAINS)
		}
		cert.Domains = userDomains
	} else {
		cert.Domains = fn.GetDomains(define.DEFAULT_DOMAINS)
	}

	k8s := strings.TrimSpace(forK8S)
	if k8s == "" {
		cert.ForK8S = define.DEFAULT_FORK8S
	} else {
		k8s = strings.ToUpper(k8s)
		if k8s == "ON" || k8s == "1" || k8s == "TRUE" {
			cert.ForK8S = "ON"
		} else {
			cert.ForK8S = "OFF"
		}
	}

	user = strings.TrimSpace(user)
	uid = strings.TrimSpace(uid)
	gid = strings.TrimSpace(gid)
	if user != "" && uid != "" && gid != "" {
		cert.OwnUser = user
		cert.OwnUID = uid
		cert.OwnGID = gid
	}

	return cert
}
