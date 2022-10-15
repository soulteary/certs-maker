package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	prepare := filepath.Join(".", "ssl")
	os.MkdirAll(prepare, os.ModePerm)
}

type CERT struct {
	Country            string
	State              string
	Locality           string
	Organization       string
	OrganizationalUnit string
	CommonName         string
	Domians            []string
}

const (
	DEFAULT_COUNTRY             = "CN"                               // Country Name
	DEFAULT_STATE               = "BJ"                               // State Or Province Name
	DEFAULT_LOCALITY            = "HD"                               // Locality Name
	DEFAULT_ORGANIZATION        = "Lab"                              // Organization Name
	DEFAULT_ORGANIZATIONAL_UNIT = "Dev"                              // Organizational Unit Name
	DEFAULT_COMMON_NAME         = "Hello World"                      // Common Name
	DEFAULT_DOMAINS             = "lab.com,*.lab.com,*.data.lab.com" // Domians
)

func verifyCountry(input string) bool {
	var re = regexp.MustCompile(`^\w{2}$`)
	ret := re.FindString(input)
	return ret != ""
}

func getDomains(input string) (result []string) {
	var re = regexp.MustCompile(`^([\.\w\*]+(\,)?){1,}$`)
	if len(re.FindAllString(input, -1)) > 0 {
		domains := strings.Split(input, ",")
		for _, domain := range domains {
			s := strings.TrimSpace(domain)
			if len(s) > 0 {
				result = append(result, strings.ToLower(domain))
			}
		}
	}
	return result
}

func getRootDomain(input string) string {
	var re = regexp.MustCompile(`([\.\w]+){1,2}$`)
	file := strings.TrimLeft(re.FindString(input), ".")
	if file == "" {
		return "cert"
	}
	return file
}

func createCertConfig(country string, state string, locality string, organization string, organizationalUnit string, commonName string, domains string) (cert CERT) {
	country = strings.TrimSpace(country)
	if len(country) > 0 {
		if verifyCountry(country) {
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
		userDomains := getDomains(domainsInput)
		if len(userDomains) == 0 {
			userDomains = getDomains(DEFAULT_DOMAINS)
			fmt.Println("wrong domains, set to default value:", DEFAULT_DOMAINS)
		}
		cert.Domians = userDomains
	} else {
		cert.Domians = getDomains(DEFAULT_DOMAINS)
	}
	return cert
}

func parseEnvInputs() (cert CERT) {
	country := os.Getenv("CERT_C")
	state := os.Getenv("CERT_ST")
	locality := os.Getenv("CERT_L")
	organization := os.Getenv("CERT_O")
	organizationalUnit := os.Getenv("CERT_OU")
	commonName := os.Getenv("CERT_CN")
	domains := os.Getenv("CERT_DNS")

	return createCertConfig(country, state, locality, organization, organizationalUnit, commonName, domains)
}

func parseCliInputs() (cert CERT) {
	var country string
	flag.StringVar(&country, "CERT_C", DEFAULT_COUNTRY, "Country Name")

	var state string
	flag.StringVar(&state, "CERT_ST", DEFAULT_STATE, "State Or Province Name")

	var locality string
	flag.StringVar(&locality, "CERT_L", DEFAULT_LOCALITY, "Locality Name")

	var organization string
	flag.StringVar(&organization, "CERT_O", DEFAULT_ORGANIZATION, "Organization Name")

	var organizationalUnit string
	flag.StringVar(&organizationalUnit, "CERT_OU", DEFAULT_ORGANIZATIONAL_UNIT, "Organizational Unit Name")

	var commonName string
	flag.StringVar(&commonName, "CERT_CN", DEFAULT_COMMON_NAME, "Common Name")

	var domains string
	flag.StringVar(&domains, "CERT_DNS", DEFAULT_DOMAINS, "Domians")

	flag.Parse()

	return createCertConfig(country, state, locality, organization, organizationalUnit, commonName, domains)
}

func mergeUserInputs() CERT {
	base := parseEnvInputs()
	cli := parseCliInputs()

	if cli.Country != DEFAULT_COUNTRY {
		base.Country = cli.Country
	}
	if cli.State != DEFAULT_STATE {
		base.State = cli.State
	}
	if cli.Locality != DEFAULT_LOCALITY {
		base.Locality = cli.Locality
	}
	if cli.Organization != DEFAULT_ORGANIZATION {
		base.Organization = cli.Organization
	}
	if cli.OrganizationalUnit != DEFAULT_ORGANIZATIONAL_UNIT {
		base.OrganizationalUnit = cli.OrganizationalUnit
	}
	if cli.CommonName != DEFAULT_COMMON_NAME {
		base.CommonName = cli.CommonName
	}
	if !reflect.DeepEqual(cli.Domians, base.Domians) {
		base.Domians = cli.Domians
	}
	return base
}

func generateConfFile(cert CERT) {

	const certBase = `
[req]
prompt                  = no
default_bits            = 4096
default_md              = sha256
encrypt_key             = no
string_mask             = utf8only

distinguished_name      = cert_distinguished_name
req_extensions          = req_x509v3_extensions
x509_extensions         = req_x509v3_extensions	
`

	certInfo := strings.Join(
		[]string{
			"[cert_distinguished_name]",
			"C  = " + cert.Country,
			"ST = " + cert.State,
			"L  = " + cert.Locality,
			"O  = " + cert.Organization,
			"OU = " + cert.OrganizationalUnit,
			"CN = " + cert.CommonName,
		}, "\n",
	)

	const certExt = `
[req_x509v3_extensions]
basicConstraints        = critical,CA:true
subjectKeyIdentifier    = hash
keyUsage                = critical,digitalSignature,keyCertSign,cRLSign
extendedKeyUsage        = critical,serverAuth
subjectAltName          = @alt_names
`

	domains := []string{"[alt_names]"}
	for idx, domain := range cert.Domians {
		id := strconv.Itoa(idx + 1)
		domains = append(domains, "DNS."+id+" = "+domain)
	}
	certDomains := strings.Join(domains, "\n")

	fileName := getRootDomain(cert.Domians[0])
	os.WriteFile("./ssl/"+fileName+".conf", []byte(certBase+"\n"+certInfo+"\n"+certExt+"\n"+certDomains), 0644)

	generateScript := "openssl req -x509 -newkey rsa:2048 -keyout ssl/${fileName}.key -out ssl/${fileName}.crt -days 3600 -nodes -config ssl/${fileName}.conf"
	os.WriteFile("./generate.sh", []byte(strings.ReplaceAll(generateScript, "${fileName}", fileName)), 0644)
}

func main() {
	config := mergeUserInputs()
	generateConfFile(config)
}
