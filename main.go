package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
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
	ForK8S             string
	OwnUser            string
	OwnUID             string
	OwnGID             string
}

const (
	DEFAULT_COUNTRY             = "CN"                               // Country Name
	DEFAULT_STATE               = "BJ"                               // State Or Province Name
	DEFAULT_LOCALITY            = "HD"                               // Locality Name
	DEFAULT_ORGANIZATION        = "Lab"                              // Organization Name
	DEFAULT_ORGANIZATIONAL_UNIT = "Dev"                              // Organizational Unit Name
	DEFAULT_COMMON_NAME         = "Hello World"                      // Common Name
	DEFAULT_DOMAINS             = "lab.com,*.lab.com,*.data.lab.com" // Domians
	DEFAULT_FORK8S              = "OFF"                              // Certs For K8S
)

var Version = "dev"

func verifyCountry(input string) bool {
	var re = regexp.MustCompile(`^\w{2}$`)
	ret := re.FindString(input)
	return ret != ""
}

func getDomains(input string) (result []string) {
	var re = regexp.MustCompile(`^([\.\w\*\-\_]+(\,)?){1,}$`)
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
	var re = regexp.MustCompile(`([\.\w\-\_]+){1,2}$`)
	file := strings.TrimLeft(re.FindString(input), ".")
	if file == "" {
		return "cert"
	}
	return file
}

func createCertConfig(country string, state string, locality string, organization string, organizationalUnit string, commonName string, domains string, forK8S string, user string, uid string, gid string) (cert CERT) {
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

func parseEnvInputs() (cert CERT) {
	country := os.Getenv("CERT_C")
	state := os.Getenv("CERT_ST")
	locality := os.Getenv("CERT_L")
	organization := os.Getenv("CERT_O")
	organizationalUnit := os.Getenv("CERT_OU")
	commonName := os.Getenv("CERT_CN")
	domains := os.Getenv("CERT_DNS")
	forK8S := os.Getenv("FOR_K8S")
	user := os.Getenv("USER")
	uid := os.Getenv("UID")
	gid := os.Getenv("GID")

	return createCertConfig(country, state, locality, organization, organizationalUnit, commonName, domains, forK8S, user, uid, gid)
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

	var forK8S string
	flag.StringVar(&forK8S, "FOR_K8S", DEFAULT_FORK8S, "FOR K8S")

	var user string
	flag.StringVar(&user, "USER", "", "File Owner Username")

	var uid string
	flag.StringVar(&uid, "UID", "", "File Owner UID")

	var gid string
	flag.StringVar(&gid, "GID", "", "File Owner GID")

	flag.Parse()

	return createCertConfig(country, state, locality, organization, organizationalUnit, commonName, domains, forK8S, user, uid, gid)
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
	if !reflect.DeepEqual(cli.Domians, getDomains(DEFAULT_DOMAINS)) {
		base.Domians = cli.Domians
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

func uniq(s []string) []string {
	k := make(map[string]bool)
	l := []string{}
	for _, i := range s {
		if _, v := k[i]; !v {
			k[i] = true
			l = append(l, i)
		}
	}
	return l
}

func generateConfFile(cert CERT) string {

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

	const certExtForK8S = `
[req_x509v3_extensions]
basicConstraints = CA:FALSE
nsCertType = server
nsComment = "OpenSSL Generated Server Certificate"
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid,issuer:always
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
`

	if cert.ForK8S == "ON" {
		cert.Domians = append(cert.Domians, "*")
		cert.Domians = append(cert.Domians, "localhost")
		cert.Domians = uniq(cert.Domians)
	}

	domains := []string{"[alt_names]"}
	for idx, domain := range cert.Domians {
		id := strconv.Itoa(idx + 1)
		domains = append(domains, "DNS."+id+" = "+domain)
	}
	certDomains := strings.Join(domains, "\n")

	fileName := getRootDomain(cert.Domians[0])
	if cert.ForK8S == "OFF" {
		os.WriteFile("./ssl/"+fileName+".conf", []byte(certBase+"\n"+certInfo+"\n"+certExt+"\n"+certDomains), 0644)
	} else {
		fileName = fileName + ".k8s"
		os.WriteFile("./ssl/"+fileName+".conf", []byte(certBase+"\n"+certInfo+"\n"+certExtForK8S+"\n"+certDomains), 0644)
	}

	scriptTpl := "openssl req -x509 -newkey rsa:2048 -keyout ./ssl/${fileName}.key -out ./ssl/${fileName}.crt -days 3650 -nodes -config ./ssl/${fileName}.conf"
	return strings.ReplaceAll(scriptTpl, "${fileName}", fileName)
}

func tryAdjustPermissions(cert CERT) {
	if cert.OwnUser == "" ||
		cert.OwnUID == "" ||
		cert.OwnGID == "" {
		return
	}
	execute(`addgroup -g ` + cert.OwnGID + ` ` + cert.OwnUser)
	execute(`adduser -g "" -G ` + cert.OwnUser + ` -H -D -u ` + cert.OwnUID + ` ` + cert.OwnUser)
	execute(`chown -R ` + cert.OwnUser + `:` + cert.OwnUser + ` ./ssl`)
	execute(`chmod -R a+r ./ssl`)
}

func execute(command string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(stdout.String())
	}
}

func main() {
	fmt.Printf("running soulteary/certs-maker %s\n", Version)
	config := mergeUserInputs()
	shell := generateConfFile(config)
	execute(shell)
	tryAdjustPermissions(config)
}
