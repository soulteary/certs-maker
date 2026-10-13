package generator_test

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/soulteary/certs-maker/internal/define"
	"github.com/soulteary/certs-maker/internal/generator"
)

func TestCertificateProfiles(t *testing.T) {
	if _, err := exec.LookPath("openssl"); err != nil {
		t.Skip("OpenSSL is required for certificate generation")
	}

	for _, test := range []struct {
		name       string
		forK8s     bool
		forFirefox bool
	}{
		{name: "default"},
		{name: "k8s", forK8s: true},
		{name: "firefox", forFirefox: true},
		{name: "k8s_takes_precedence", forK8s: true, forFirefox: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			configureCertificateProfile(t, test.forK8s, test.forFirefox)
			generator.MakeCerts()

			name := "example.test"
			if test.forK8s {
				name += ".k8s"
			}
			base := filepath.Join(define.APP_OUTPUT_DIR, name)
			cert := readProfileCertificate(t, base+".pem.crt")
			isFirefox := test.forFirefox && !test.forK8s
			if !cert.BasicConstraintsValid || cert.IsCA != (!test.forK8s && !isFirefox) {
				t.Fatalf("unexpected certificate CA constraints: valid=%v, CA=%v", cert.BasicConstraintsValid, cert.IsCA)
			}
			if len(cert.ExtKeyUsage) != 1 || cert.ExtKeyUsage[0] != x509.ExtKeyUsageServerAuth {
				t.Fatalf("expected serverAuth usage, got %v", cert.ExtKeyUsage)
			}
			for _, host := range []string{"example.test", "service.example.test", "127.0.0.1", "::1"} {
				if err := cert.VerifyHostname(host); err != nil {
					t.Errorf("certificate does not cover %s: %v", host, err)
				}
			}
			if err := cert.VerifyHostname("unrelated.test"); err == nil {
				t.Error("certificate unexpectedly covers an unrelated hostname")
			}
			if test.forK8s {
				if err := cert.VerifyHostname("localhost"); err != nil {
					t.Errorf("K8s certificate does not cover localhost: %v", err)
				}
			}
			der, err := os.ReadFile(base + ".der.crt")
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(der, cert.Raw) {
				t.Error("DER certificate does not match the PEM certificate")
			}

			if !isFirefox {
				if err := cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature); err != nil {
					t.Errorf("expected a self-signed certificate: %v", err)
				}
				if _, err := os.Stat(base + ".rootCA.pem"); !os.IsNotExist(err) {
					t.Errorf("unexpected Firefox root certificate: %v", err)
				}
				return
			}

			rootPath := base + ".rootCA.pem"
			root := readProfileCertificate(t, rootPath)
			if !root.BasicConstraintsValid || !root.IsCA || root.KeyUsage != x509.KeyUsageCertSign|x509.KeyUsageCRLSign {
				t.Fatalf("unexpected root CA constraints: valid=%v, CA=%v, usage=%v", root.BasicConstraintsValid, root.IsCA, root.KeyUsage)
			}
			if len(root.DNSNames) != 0 || len(root.IPAddresses) != 0 || len(root.ExtKeyUsage) != 0 {
				t.Error("root CA unexpectedly includes server certificate extensions")
			}
			if err := root.CheckSignatureFrom(root); err != nil {
				t.Errorf("root CA is not self-signed: %v", err)
			}
			if cert.KeyUsage != x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment {
				t.Errorf("unexpected leaf key usage: %v", cert.KeyUsage)
			}
			if err := cert.CheckSignatureFrom(root); err != nil {
				t.Errorf("leaf certificate is not signed by the root CA: %v", err)
			}
			roots := x509.NewCertPool()
			roots.AddCert(root)
			for _, host := range []string{"example.test", "service.example.test", "127.0.0.1", "::1"} {
				if _, err := cert.Verify(x509.VerifyOptions{Roots: roots, DNSName: host}); err != nil {
					t.Errorf("Go cannot verify the TLS chain for %s: %v", host, err)
				}
			}
			if output, err := exec.Command("openssl", "verify", "-CAfile", rootPath, "-purpose", "sslserver", base+".pem.crt").CombinedOutput(); err != nil {
				t.Errorf("OpenSSL cannot verify the TLS certificate chain: %v\n%s", err, output)
			}
		})
	}
}

func readProfileCertificate(t *testing.T, path string) *x509.Certificate {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		t.Fatalf("invalid PEM certificate in %s", path)
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	return cert
}

func configureCertificateProfile(t *testing.T, forK8s, forFirefox bool) {
	t.Helper()
	country, state, locality := define.CERT_COUNTRY, define.CERT_STATE, define.CERT_LOCALITY
	organization, unit, commonName := define.CERT_ORGANIZATION, define.CERT_ORGANIZATIONAL_UNIT, define.CERT_COMMON_NAME
	domains := append([]string(nil), define.CERT_DOMAINS...)
	outputDir, fileName, expireDays := define.APP_OUTPUT_DIR, define.CUSTOM_FILE_NAME, define.CERT_EXPIRE_DAYS
	k8s, firefox := define.APP_FOR_K8S, define.APP_FOR_FIREFOX
	t.Cleanup(func() {
		define.CERT_COUNTRY, define.CERT_STATE, define.CERT_LOCALITY = country, state, locality
		define.CERT_ORGANIZATION, define.CERT_ORGANIZATIONAL_UNIT, define.CERT_COMMON_NAME = organization, unit, commonName
		define.CERT_DOMAINS = domains
		define.APP_OUTPUT_DIR, define.CUSTOM_FILE_NAME, define.CERT_EXPIRE_DAYS = outputDir, fileName, expireDays
		define.APP_FOR_K8S, define.APP_FOR_FIREFOX = k8s, firefox
	})
	define.CERT_COUNTRY, define.CERT_STATE, define.CERT_LOCALITY = define.DEFAULT_COUNTRY, define.DEFAULT_STATE, define.DEFAULT_LOCALITY
	define.CERT_ORGANIZATION, define.CERT_ORGANIZATIONAL_UNIT, define.CERT_COMMON_NAME = define.DEFAULT_ORGANIZATION, define.DEFAULT_ORGANIZATIONAL_UNIT, define.DEFAULT_COMMON_NAME
	define.CERT_DOMAINS = []string{"example.test", "*.example.test", "127.0.0.1", "::1"}
	define.APP_OUTPUT_DIR, define.CUSTOM_FILE_NAME, define.CERT_EXPIRE_DAYS = t.TempDir(), "", define.DEFAULT_EXPIRE_DAYS
	define.APP_FOR_K8S, define.APP_FOR_FIREFOX = forK8s, forFirefox
}
