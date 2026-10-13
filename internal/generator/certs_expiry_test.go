/**
 * Copyright (c) 2021-2026 Su Yang (soulteary)
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package generator_test

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/soulteary/certs-maker/internal/define"
	"github.com/soulteary/certs-maker/internal/generator"
)

func TestMakeCertsCertificateExpiry(t *testing.T) {
	if _, err := exec.LookPath("openssl"); err != nil {
		t.Skip("OpenSSL is required for certificate generation")
	}

	for _, tc := range []struct {
		name       string
		expireDays string
		days       int
		k8s        bool
		firefox    bool
	}{
		{name: "custom", expireDays: "7", days: 7},
		{name: "default", days: 3650},
		{name: "blank", expireDays: " \t", days: 3650},
		{name: "k8s", expireDays: "3", days: 3, k8s: true},
		{name: "firefox", expireDays: "2", days: 2, firefox: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := prepareCertificateTest(t)
			define.CERT_EXPIRE_DAYS = tc.expireDays
			define.APP_FOR_K8S = tc.k8s
			define.APP_FOR_FIREFOX = tc.firefox
			if err := generator.MakeCerts(); err != nil {
				t.Fatalf("generate certificate: %v", err)
			}

			name := "expiry.example.com"
			if tc.k8s {
				name += ".k8s"
			}
			assertCertificateExpiry(t, filepath.Join(dir, name+".pem.crt"), tc.days)
			if tc.firefox {
				assertCertificateExpiry(t, filepath.Join(dir, name+".rootCA.pem"), tc.days)
			}
		})
	}
}

func TestMakeCertsRejectsInvalidExpiry(t *testing.T) {
	for _, expireDays := range []string{"0", "-1", "1.5", "invalid", "2147483648", "7; touch injected"} {
		t.Run(expireDays, func(t *testing.T) {
			dir := prepareCertificateTest(t)
			define.CERT_EXPIRE_DAYS = expireDays
			if err := generator.Generate(); err == nil {
				t.Fatal("invalid expiry must return an error")
			}

			files, err := os.ReadDir(dir)
			if err != nil {
				t.Fatal(err)
			}
			if len(files) != 0 {
				t.Fatalf("invalid expiry must not generate files: %v", files)
			}
		})
	}
}

func TestMakeCertsReturnsWriteError(t *testing.T) {
	dir := prepareCertificateTest(t)
	define.CERT_EXPIRE_DAYS = "7"
	define.APP_OUTPUT_DIR = filepath.Join(dir, "missing")
	if err := generator.MakeCerts(); err == nil || !strings.Contains(err.Error(), "write certificate configuration") {
		t.Fatalf("expected certificate configuration write error; got %v", err)
	}
}

func TestMakeCertsStopsOnOpenSSLFailure(t *testing.T) {
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("sh is required for certificate generation")
	}
	dir := prepareCertificateTest(t)
	define.CERT_EXPIRE_DAYS = "7"
	binDir := t.TempDir()
	if err := os.Symlink(sh, filepath.Join(binDir, "sh")); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir)
	if err := generator.Generate(); err == nil || !strings.Contains(err.Error(), "generate certificate") {
		t.Fatalf("expected OpenSSL generation error; got %v", err)
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0].Name() != "expiry.example.com.conf" {
		t.Fatalf("failed generation must stop before DER conversion: %v", files)
	}
}

func prepareCertificateTest(t *testing.T) string {
	t.Helper()
	country, state, locality := define.CERT_COUNTRY, define.CERT_STATE, define.CERT_LOCALITY
	organization, unit, commonName := define.CERT_ORGANIZATION, define.CERT_ORGANIZATIONAL_UNIT, define.CERT_COMMON_NAME
	domains := define.CERT_DOMAINS
	k8s, firefox := define.APP_FOR_K8S, define.APP_FOR_FIREFOX
	dir, customName, expiry := define.APP_OUTPUT_DIR, define.CUSTOM_FILE_NAME, define.CERT_EXPIRE_DAYS
	user, uid, gid := define.APP_USER, define.APP_UID, define.APP_GID
	t.Cleanup(func() {
		define.CERT_COUNTRY, define.CERT_STATE, define.CERT_LOCALITY = country, state, locality
		define.CERT_ORGANIZATION, define.CERT_ORGANIZATIONAL_UNIT, define.CERT_COMMON_NAME = organization, unit, commonName
		define.CERT_DOMAINS = domains
		define.APP_FOR_K8S, define.APP_FOR_FIREFOX = k8s, firefox
		define.APP_OUTPUT_DIR, define.CUSTOM_FILE_NAME, define.CERT_EXPIRE_DAYS = dir, customName, expiry
		define.APP_USER, define.APP_UID, define.APP_GID = user, uid, gid
	})
	define.CERT_COUNTRY = define.DEFAULT_COUNTRY
	define.CERT_STATE = define.DEFAULT_STATE
	define.CERT_LOCALITY = define.DEFAULT_LOCALITY
	define.CERT_ORGANIZATION = define.DEFAULT_ORGANIZATION
	define.CERT_ORGANIZATIONAL_UNIT = define.DEFAULT_ORGANIZATIONAL_UNIT
	define.CERT_COMMON_NAME = "expiry.example.com"
	define.CERT_DOMAINS = []string{"expiry.example.com"}
	define.APP_FOR_K8S, define.APP_FOR_FIREFOX = false, false
	define.CUSTOM_FILE_NAME = ""
	define.CERT_EXPIRE_DAYS = define.DEFAULT_EXPIRE_DAYS
	define.APP_USER, define.APP_UID, define.APP_GID = "", "", ""
	define.APP_OUTPUT_DIR = t.TempDir()
	return define.APP_OUTPUT_DIR
}

func assertCertificateExpiry(t *testing.T, path string, days int) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated certificate: %v", err)
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		t.Fatalf("%s is not a PEM certificate", path)
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse generated certificate: %v", err)
	}
	if got, want := cert.NotAfter.Sub(cert.NotBefore), time.Duration(days)*24*time.Hour; got != want {
		t.Fatalf("%s validity is %v; want %v", path, got, want)
	}
	if delta := time.Since(cert.NotBefore); delta < -time.Second || delta > time.Minute {
		t.Fatalf("%s was not issued in this test: NotBefore=%v", path, cert.NotBefore)
	}
}
