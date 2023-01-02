package cmd

import (
	"flag"
	"strings"

	"github.com/soulteary/certs-maker/internal/define"
)

func ParseFlags() (appFlags AppFlags) {
	flag.StringVar(&appFlags.Country, ENV_KEY_COUNTRY, define.DEFAULT_COUNTRY, CLI_DESC_COUNTRY)
	flag.StringVar(&appFlags.State, ENV_KEY_STATE, define.DEFAULT_STATE, CLI_DESC_STATE)
	flag.StringVar(&appFlags.Locality, ENV_KEY_LOCALITY, define.DEFAULT_LOCALITY, CLI_DESC_LOCALITY)
	flag.StringVar(&appFlags.Organization, ENV_KEY_ORGANIZATION, define.DEFAULT_ORGANIZATION, CLI_DESC_ORGANIZATION)
	flag.StringVar(&appFlags.OrganizationalUnit, ENV_KEY_ORGANIZATION_UNIT, define.DEFAULT_ORGANIZATIONAL_UNIT, CLI_DESC_ORGANIZATION_UNIT)
	flag.StringVar(&appFlags.CommonName, ENV_KEY_COMMON_NAME, define.DEFAULT_COMMON_NAME, CLI_DESC_COMMON_NAME)
	flag.StringVar(&appFlags.Domains, ENV_KEY_DOMAINS, define.DEFAULT_DOMAINS, CLI_DESC_DOMAINS)

	flag.BoolVar(&appFlags.ForK8s, ENV_KEY_FOR_K8S, define.DEFAULT_FOR_K8S, CLI_DESC_FOR_K8S)
	flag.StringVar(&appFlags.User, ENV_KEY_USER, define.DEFAULT_USER, CLI_DESC_USER)
	flag.StringVar(&appFlags.UID, ENV_KEY_UID, define.DEFAULT_UID, CLI_DESC_UID)
	flag.StringVar(&appFlags.GID, ENV_KEY_GID, define.APP_GID, CLI_DESC_GID)

	flag.Parse()
	return appFlags
}

func ApplyFlags() {
	args := ParseFlags()

	define.CERT_COUNTRY = UpdateCountryOption(ENV_KEY_COUNTRY, args.Country, define.DEFAULT_COUNTRY)
	define.CERT_STATE = UpdateStringOption(ENV_KEY_STATE, args.State, define.DEFAULT_STATE)
	define.CERT_LOCALITY = strings.ToUpper(UpdateStringOption(ENV_KEY_LOCALITY, args.Locality, define.DEFAULT_LOCALITY))
	define.CERT_ORGANIZATION = UpdateStringOption(ENV_KEY_ORGANIZATION, args.Organization, define.DEFAULT_ORGANIZATION)
	define.CERT_ORGANIZATIONAL_UNIT = UpdateStringOption(ENV_KEY_ORGANIZATION_UNIT, args.OrganizationalUnit, define.DEFAULT_ORGANIZATIONAL_UNIT)
	define.CERT_COMMON_NAME = UpdateStringOption(ENV_KEY_COMMON_NAME, args.CommonName, define.DEFAULT_COMMON_NAME)
	define.CERT_DOMAINS = UpdateDomainOption(ENV_KEY_DOMAINS, args.Domains, define.DEFAULT_DOMAINS)
	define.APP_FOR_K8S = UpdateBoolOption(ENV_KEY_FOR_K8S, args.ForK8s, define.DEFAULT_FOR_K8S)

	user := UpdateStringOption(ENV_KEY_USER, args.User, define.DEFAULT_USER)
	uid := UpdateStringOption(ENV_KEY_UID, args.UID, define.DEFAULT_UID)
	gid := UpdateStringOption(ENV_KEY_GID, args.GID, define.DEFAULT_GID)

	if user != "" && uid != "" && gid != "" {
		define.APP_USER = user
		define.APP_UID = uid
		define.APP_GID = gid
	}
}