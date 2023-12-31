package cmd

import (
	"fmt"

	"github.com/soulteary/certs-maker/internal/define"
)

const (
	ENV_KEY_COUNTRY           = "CERT_C"
	ENV_KEY_STATE             = "CERT_ST"
	ENV_KEY_LOCALITY          = "CERT_L"
	ENV_KEY_ORGANIZATION      = "CERT_O"
	ENV_KEY_ORGANIZATION_UNIT = "CERT_OU"
	ENV_KEY_COMMON_NAME       = "CERT_CN"
	ENV_KEY_DOMAINS           = "CERT_DNS"

	ENV_KEY_FOR_K8S     = "FOR_K8S"
	ENV_KEY_FOR_FIREFOX = "FOR_FIREFOX"

	ENV_KEY_USER       = "USER"
	ENV_KEY_UID        = "UID"
	ENV_KEY_GID        = "GID"
	ENV_KEY_OUTPUT_DIR = "DIR"
)

var (
	CLI_DESC_COUNTRY           = fmt.Sprintf("Country Name, env: `%s`, default: `%s`", ENV_KEY_COUNTRY, define.DEFAULT_COUNTRY)
	CLI_DESC_STATE             = fmt.Sprintf("State Or Province Name, env: `%s`, default: `%s`", ENV_KEY_STATE, define.DEFAULT_STATE)
	CLI_DESC_LOCALITY          = fmt.Sprintf("Locality Name, env: `%s`, default: `%s`", ENV_KEY_LOCALITY, define.DEFAULT_LOCALITY)
	CLI_DESC_ORGANIZATION      = fmt.Sprintf("Organization Name, env: `%s`, default: `%s`", ENV_KEY_ORGANIZATION, define.DEFAULT_ORGANIZATION)
	CLI_DESC_ORGANIZATION_UNIT = fmt.Sprintf("Organizational Unit Name, env: `%s`, default: `%s`", ENV_KEY_ORGANIZATION_UNIT, define.DEFAULT_ORGANIZATIONAL_UNIT)
	CLI_DESC_COMMON_NAME       = fmt.Sprintf("Common Name, env: `%s`, default: `%s`", ENV_KEY_COMMON_NAME, define.DEFAULT_COMMON_NAME)
	CLI_DESC_DOMAINS           = fmt.Sprintf("Domains, env: `%s`, default: `%s`", ENV_KEY_DOMAINS, define.DEFAULT_DOMAINS)
	CLI_DESC_FOR_K8S           = fmt.Sprintf("Issue for K8s, env: `%s`, default: `%v`", ENV_KEY_FOR_K8S, define.DEFAULT_FOR_K8S)
	CLI_DESC_USER              = fmt.Sprintf("File Owner User, env: `%s`, default: `%s`", ENV_KEY_USER, define.DEFAULT_USER)
	CLI_DESC_UID               = fmt.Sprintf("File Owner UID, env: `%s`, default: `%s`", ENV_KEY_UID, define.DEFAULT_UID)
	CLI_DESC_GID               = fmt.Sprintf("File Owner GID, env: `%s`, default: `%s`", ENV_KEY_GID, define.DEFAULT_GID)
	CLI_DESC_OUTPUT_DIR        = fmt.Sprintf("Certs Dir, env: `%s`, default: `%s`", ENV_KEY_OUTPUT_DIR, define.DEFAULT_DIR)
)

type AppFlags struct {
	Country            string
	State              string
	Locality           string
	Organization       string
	OrganizationalUnit string
	CommonName         string
	Domains            string

	ForK8s     string
	ForFirefox string

	User      string
	UID       string
	GID       string
	OutputDir string
}
