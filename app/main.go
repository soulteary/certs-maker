package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
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
	return strings.TrimLeft(re.FindString(input), ".")
}

func parseUserInputs() (cert CERT) {
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

	country = strings.TrimSpace(country)
	if len(country) > 0 {
		if verifyCountry(country) {
			cert.Country = strings.ToUpper(country)
		} else {
			fmt.Println("wrong country name, set to default value:", DEFAULT_COUNTRY)
		}
	}

	state = strings.TrimSpace(state)
	if len(state) > 0 {
		cert.State = strings.ToUpper(state)
	}

	locality = strings.TrimSpace(locality)
	if len(locality) > 0 {
		cert.Locality = strings.ToUpper(locality)
	}

	organization = strings.TrimSpace(organization)
	if len(organization) > 0 {
		cert.Organization = organization
	}

	organizationalUnit = strings.TrimSpace(organizationalUnit)
	if len(organization) > 0 {
		cert.OrganizationalUnit = organizationalUnit
	}

	commonName = strings.TrimSpace(commonName)
	if len(commonName) > 0 {
		cert.CommonName = commonName
	}

	userDomains := getDomains(strings.TrimSpace(domains))
	if len(userDomains) == 0 {
		userDomains = getDomains(DEFAULT_DOMAINS)
		fmt.Println("wrong domains, set to default value:", DEFAULT_DOMAINS)
	}
	cert.Domians = userDomains
	return cert
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
	config := parseUserInputs()
	generateConfFile(config)
}
