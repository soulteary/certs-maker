package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"cert-maker/internal/cmd"
	"cert-maker/internal/define"
	"cert-maker/internal/fn"
	"cert-maker/internal/version"
)

func init() {
	prepare := filepath.Join(".", "ssl")
	os.MkdirAll(prepare, os.ModePerm)
}

func generateConfFile(cert define.CERT) string {

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
		cert.Domains = append(cert.Domains, "*")
		cert.Domains = append(cert.Domains, "localhost")
		cert.Domains = fn.Uniq(cert.Domains)
	}

	domains := []string{"[alt_names]"}
	for idx, domain := range cert.Domains {
		id := strconv.Itoa(idx + 1)
		domains = append(domains, "DNS."+id+" = "+domain)
	}
	certDomains := strings.Join(domains, "\n")

	fileName := fn.GetRootDomain(cert.Domains[0])
	if cert.ForK8S == "OFF" {
		os.WriteFile("./ssl/"+fileName+".conf", []byte(certBase+"\n"+certInfo+"\n"+certExt+"\n"+certDomains), 0644)
	} else {
		fileName = fileName + ".k8s"
		os.WriteFile("./ssl/"+fileName+".conf", []byte(certBase+"\n"+certInfo+"\n"+certExtForK8S+"\n"+certDomains), 0644)
	}

	scriptTpl := "openssl req -x509 -newkey rsa:2048 -keyout ./ssl/${fileName}.key -out ./ssl/${fileName}.crt -days 3650 -nodes -config ./ssl/${fileName}.conf"
	return strings.ReplaceAll(scriptTpl, "${fileName}", fileName)
}

func tryAdjustPermissions(cert define.CERT) {
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
	fmt.Printf("running soulteary/certs-maker %s\n", version.Version)
	config := cmd.MergeUserInputs()
	shell := generateConfFile(config)
	execute(shell)
	tryAdjustPermissions(config)
}
