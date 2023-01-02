package cmd

import (
	"cert-maker/internal/define"
	"flag"
)

func ParseFlags() (cert define.CERT) {
	var country string
	flag.StringVar(&country, ENV_KEY_COUNTRY, DEFAULT_COUNTRY, "Country Name")

	var state string
	flag.StringVar(&state, ENV_KEY_STATE, DEFAULT_STATE, "State Or Province Name")

	var locality string
	flag.StringVar(&locality, ENV_KEY_LOCALITY, DEFAULT_LOCALITY, "Locality Name")

	var organization string
	flag.StringVar(&organization, ENV_KEY_ORGANIZATION, DEFAULT_ORGANIZATION, "Organization Name")

	var organizationalUnit string
	flag.StringVar(&organizationalUnit, ENV_KEY_ORGANIZATION_UNIT, DEFAULT_ORGANIZATIONAL_UNIT, "Organizational Unit Name")

	var commonName string
	flag.StringVar(&commonName, ENV_KEY_COMMON_NAME, DEFAULT_COMMON_NAME, "Common Name")

	var domains string
	flag.StringVar(&domains, ENV_KEY_DOMAINS, DEFAULT_DOMAINS, "Domains")

	var forK8S string
	flag.StringVar(&forK8S, ENV_KEY_FOR_K8S, DEFAULT_FORK8S, "FOR K8S")

	var user string
	flag.StringVar(&user, ENV_KEY_USER, "", "File Owner Username")

	var uid string
	flag.StringVar(&uid, ENV_KEY_UID, "", "File Owner UID")

	var gid string
	flag.StringVar(&gid, ENV_KEY_GID, "", "File Owner GID")

	flag.Parse()

	return createCertConfig(country, state, locality, organization, organizationalUnit, commonName, domains, forK8S, user, uid, gid)
}
