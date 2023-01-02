package cmd

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"cert-maker/internal/define"
	"cert-maker/internal/fn"
)

func parseEnvInputs() (cert define.CERT) {
	country := os.Getenv(ENV_KEY_COUNTRY)
	state := os.Getenv(ENV_KEY_STATE)
	locality := os.Getenv(ENV_KEY_LOCALITY)
	organization := os.Getenv(ENV_KEY_ORGANIZATION)
	organizationalUnit := os.Getenv(ENV_KEY_ORGANIZATION_UNIT)
	commonName := os.Getenv(ENV_KEY_COMMON_NAME)
	domains := os.Getenv(ENV_KEY_DOMAINS)
	forK8S := os.Getenv(ENV_KEY_FOR_K8S)
	user := os.Getenv(ENV_KEY_USER)
	uid := os.Getenv(ENV_KEY_UID)
	gid := os.Getenv(ENV_KEY_GID)

	return createCertConfig(country, state, locality, organization, organizationalUnit, commonName, domains, forK8S, user, uid, gid)
}

func MergeUserInputs() define.CERT {
	base := parseEnvInputs()
	args := ParseFlags()

	if args.Country != DEFAULT_COUNTRY {
		base.Country = args.Country
	}
	if args.State != DEFAULT_STATE {
		base.State = args.State
	}
	if args.Locality != DEFAULT_LOCALITY {
		base.Locality = args.Locality
	}
	if args.Organization != DEFAULT_ORGANIZATION {
		base.Organization = args.Organization
	}
	if args.OrganizationalUnit != DEFAULT_ORGANIZATIONAL_UNIT {
		base.OrganizationalUnit = args.OrganizationalUnit
	}
	if args.CommonName != DEFAULT_COMMON_NAME {
		base.CommonName = args.CommonName
	}
	if !reflect.DeepEqual(args.Domains, fn.GetDomains(DEFAULT_DOMAINS)) {
		base.Domains = args.Domains
	}
	if args.ForK8S != base.ForK8S {
		base.ForK8S = args.ForK8S
	}
	if args.OwnUser != "" && args.OwnUID != "" && args.OwnGID != "" {
		if args.OwnUser != base.OwnUser && args.OwnUID != base.OwnUID && args.OwnGID != base.OwnGID {
			base.OwnUser = args.OwnUser
			base.OwnUID = args.OwnUID
			base.OwnGID = args.OwnGID
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
			fmt.Println("wrong country name, set to default value:", DEFAULT_COUNTRY)
		}
	} else {
		cert.Country = DEFAULT_COUNTRY
	}

	state = strings.TrimSpace(state)
	if len(state) > 0 {
		cert.State = strings.ToUpper(state)
	} else {
		cert.State = DEFAULT_STATE
	}

	locality = strings.TrimSpace(locality)
	if len(locality) > 0 {
		cert.Locality = strings.ToUpper(locality)
	} else {
		cert.Locality = DEFAULT_LOCALITY
	}

	organization = strings.TrimSpace(organization)
	if len(organization) > 0 {
		cert.Organization = organization
	} else {
		cert.Organization = DEFAULT_ORGANIZATION
	}

	organizationalUnit = strings.TrimSpace(organizationalUnit)
	if len(organization) > 0 {
		cert.OrganizationalUnit = organizationalUnit
	} else {
		cert.OrganizationalUnit = DEFAULT_ORGANIZATIONAL_UNIT
	}

	commonName = strings.TrimSpace(commonName)
	if len(commonName) > 0 {
		cert.CommonName = commonName
	} else {
		cert.CommonName = DEFAULT_COMMON_NAME
	}

	domainsInput := strings.TrimSpace(domains)
	if len(domainsInput) > 0 {
		userDomains := fn.GetDomains(domainsInput)
		if len(userDomains) == 0 {
			userDomains = fn.GetDomains(DEFAULT_DOMAINS)
			fmt.Println("wrong domains, set to default value:", DEFAULT_DOMAINS)
		}
		cert.Domains = userDomains
	} else {
		cert.Domains = fn.GetDomains(DEFAULT_DOMAINS)
	}

	k8s := strings.TrimSpace(forK8S)
	if k8s == "" {
		cert.ForK8S = DEFAULT_FORK8S
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
