package cmd

import (
	"flag"
	"fmt"
	"os"
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

	flag.StringVar(&appFlags.ForK8s, ENV_KEY_FOR_K8S, define.DEFAULT_FOR_K8S, CLI_DESC_FOR_K8S)
	flag.StringVar(&appFlags.ForFirefox, ENV_KEY_FOR_FIREFOX, define.DEFAULT_FOR_FIREFOX, CLI_DESC_FOR_FIREFOX)

	flag.StringVar(&appFlags.User, ENV_KEY_USER, define.DEFAULT_USER, CLI_DESC_USER)
	flag.StringVar(&appFlags.UID, ENV_KEY_UID, define.DEFAULT_UID, CLI_DESC_UID)
	flag.StringVar(&appFlags.GID, ENV_KEY_GID, define.APP_GID, CLI_DESC_GID)
	flag.StringVar(&appFlags.OutputDir, ENV_KEY_OUTPUT_DIR, define.APP_OUTPUT_DIR, CLI_DESC_OUTPUT_DIR)

	flag.StringVar(&appFlags.ExpireDays, ENV_KEY_EXPIRE_DAYS, define.DEFAULT_EXPIRE_DAYS, CLI_DESC_EXPIRE_DAYS)
	flag.StringVar(&appFlags.CustomFileName, ENV_KEY_CUSTOM_FILE_NAME, define.DEFAULT_CUSTOM_FILE_NAME, CLI_DESC_CUSTOM_FILE_NAME)

	flag.Parse()
	return appFlags
}

func ApplyFlags() {
	args := ParseFlags()
	fmt.Println("Flags:")

	define.CERT_COUNTRY = UpdateCountryOption(ENV_KEY_COUNTRY, args.Country, define.DEFAULT_COUNTRY)
	fmt.Println("  - CERT_COUNTRY=", define.CERT_COUNTRY)

	define.CERT_STATE = UpdateStringOption(ENV_KEY_STATE, args.State, define.DEFAULT_STATE)
	fmt.Println("  - CERT_STATE=", define.CERT_STATE)

	define.CERT_LOCALITY = strings.ToUpper(UpdateStringOption(ENV_KEY_LOCALITY, args.Locality, define.DEFAULT_LOCALITY))
	fmt.Println("  - CERT_LOCALITY=", define.CERT_LOCALITY)

	define.CERT_ORGANIZATION = UpdateStringOption(ENV_KEY_ORGANIZATION, args.Organization, define.DEFAULT_ORGANIZATION)
	fmt.Println("  - CERT_ORGANIZATION=", define.CERT_ORGANIZATION)

	define.CERT_ORGANIZATIONAL_UNIT = UpdateStringOption(ENV_KEY_ORGANIZATION_UNIT, args.OrganizationalUnit, define.DEFAULT_ORGANIZATIONAL_UNIT)
	fmt.Println("  - CERT_ORGANIZATIONAL_UNIT=", define.CERT_ORGANIZATIONAL_UNIT)

	define.CERT_COMMON_NAME = UpdateStringOption(ENV_KEY_COMMON_NAME, args.CommonName, define.DEFAULT_COMMON_NAME)
	fmt.Println("  - CERT_COMMON_NAME=", define.CERT_COMMON_NAME)

	define.CERT_DOMAINS = UpdateDomainOption(ENV_KEY_DOMAINS, args.Domains, define.DEFAULT_DOMAINS)
	fmt.Println("  - CERT_DOMAINS=", define.CERT_DOMAINS)

	define.APP_FOR_K8S = UpdateBoolOption(ENV_KEY_FOR_K8S, args.ForK8s, define.DEFAULT_FOR_K8S)
	fmt.Println("  - APP_FOR_K8S=", define.APP_FOR_K8S)

	define.APP_FOR_FIREFOX = UpdateBoolOption(ENV_KEY_FOR_FIREFOX, args.ForFirefox, define.DEFAULT_FOR_FIREFOX)
	fmt.Println("  - APP_FOR_FIREFOX=", define.APP_FOR_FIREFOX)

	define.APP_OUTPUT_DIR = SantizeDirPath(ENV_KEY_OUTPUT_DIR, args.OutputDir, define.DEFAULT_DIR)
	os.MkdirAll(define.APP_OUTPUT_DIR, os.ModePerm)
	fmt.Println("  - APP_OUTPUT_DIR=", define.APP_OUTPUT_DIR)
	define.CUSTOM_FILE_NAME = UpdateStringOption(ENV_KEY_CUSTOM_FILE_NAME, args.CustomFileName, define.DEFAULT_CUSTOM_FILE_NAME)
	fmt.Println("  - CUSTOM_FILE_NAME=", define.CUSTOM_FILE_NAME)

	user := UpdateStringOption(ENV_KEY_USER, args.User, define.DEFAULT_USER)
	uid := UpdateStringOption(ENV_KEY_UID, args.UID, define.DEFAULT_UID)
	gid := UpdateStringOption(ENV_KEY_GID, args.GID, define.DEFAULT_GID)

	if user != "" && uid != "" && gid != "" {
		define.APP_USER = user
		fmt.Println("  - APP_USER=", define.APP_USER)

		define.APP_UID = uid
		fmt.Println("  - APP_UID=", define.APP_UID)

		define.APP_GID = gid
		fmt.Println("  - APP_GID=", define.APP_GID)
	}

	define.CERT_EXPIRE_DAYS = UpdateStringOption(ENV_KEY_EXPIRE_DAYS, args.ExpireDays, define.DEFAULT_EXPIRE_DAYS)
	fmt.Println("  - CERT_EXPIRE_DAYS=", define.CERT_EXPIRE_DAYS)
}
