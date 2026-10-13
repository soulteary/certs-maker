/**
 * Copyright (c) 2021-2025 Su Yang (soulteary)
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package cmd_test

import (
	"flag"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/soulteary/certs-maker/internal/cmd"
	"github.com/soulteary/certs-maker/internal/define"
	"github.com/soulteary/certs-maker/internal/fn"
)

func TestApplyFlags(t *testing.T) {
	const user = "soulteary"
	const uid = "1234"
	const gid = "4321"

	os.Setenv(cmd.ENV_KEY_USER, user)
	os.Setenv(cmd.ENV_KEY_UID, uid)
	os.Setenv(cmd.ENV_KEY_GID, gid)

	cmd.ApplyFlags()

	if define.CERT_COUNTRY != define.DEFAULT_COUNTRY {
		t.Fatal("test flag parse failed")
	}

	if define.CERT_STATE != define.DEFAULT_STATE {
		t.Fatal("test flag parse failed")
	}

	if define.CERT_LOCALITY != define.DEFAULT_LOCALITY {
		t.Fatal("test flag parse failed")
	}

	if define.CERT_ORGANIZATION != define.DEFAULT_ORGANIZATION {
		t.Fatal("test flag parse failed")
	}

	if define.CERT_ORGANIZATIONAL_UNIT != define.DEFAULT_ORGANIZATIONAL_UNIT {
		t.Fatal("test flag parse failed")
	}

	if define.CERT_COMMON_NAME != define.DEFAULT_COMMON_NAME {
		t.Fatal("test flag parse failed")
	}

	if !reflect.DeepEqual(define.CERT_DOMAINS, fn.GetDomainsByString(define.DEFAULT_DOMAINS)) {
		t.Fatal("test flag parse failed")
	}

	if define.APP_FOR_K8S != fn.IsBoolString(define.DEFAULT_FOR_K8S) {
		t.Fatal("test flag parse failed")
	}

	if define.APP_USER != user {
		t.Fatal("test flag parse failed")
	}
	if define.APP_UID != uid {
		t.Fatal("test flag parse failed")
	}
	if define.APP_GID != gid {
		t.Fatal("test flag parse failed")
	}

	if define.APP_OUTPUT_DIR != define.DEFAULT_DIR {
		t.Fatal("test flag parse failed")
	}

	if define.CERT_EXPIRE_DAYS != define.DEFAULT_EXPIRE_DAYS {
		t.Fatal("test flag parse failed")
	}

	if define.CUSTOM_FILE_NAME != define.DEFAULT_CUSTOM_FILE_NAME {
		t.Fatal("test flag parse failed")
	}
}

func TestApplyFlagsPrecedence(t *testing.T) {
	originalFlags, originalArgs := flag.CommandLine, os.Args
	t.Cleanup(func() {
		flag.CommandLine, os.Args = originalFlags, originalArgs
	})
	t.Chdir(t.TempDir())

	defaults := cmd.AppFlags{
		Country:            define.DEFAULT_COUNTRY,
		State:              define.DEFAULT_STATE,
		Locality:           define.DEFAULT_LOCALITY,
		Organization:       define.DEFAULT_ORGANIZATION,
		OrganizationalUnit: define.DEFAULT_ORGANIZATIONAL_UNIT,
		CommonName:         define.DEFAULT_COMMON_NAME,
		Domains:            define.DEFAULT_DOMAINS,
		ForK8s:             "false",
		ForFirefox:         "false",
		User:               "certowner",
		UID:                "1000",
		GID:                "1000",
		OutputDir:          define.DEFAULT_DIR,
		ExpireDays:         define.DEFAULT_EXPIRE_DAYS,
		CustomFileName:     define.DEFAULT_CUSTOM_FILE_NAME,
	}
	environment := map[string]string{
		cmd.ENV_KEY_COUNTRY:           "US",
		cmd.ENV_KEY_STATE:             "Environment State",
		cmd.ENV_KEY_LOCALITY:          "environment locality",
		cmd.ENV_KEY_ORGANIZATION:      "Environment Organization",
		cmd.ENV_KEY_ORGANIZATION_UNIT: "Environment Unit",
		cmd.ENV_KEY_COMMON_NAME:       "Environment Common Name",
		cmd.ENV_KEY_DOMAINS:           "env.example.com,*.env.example.com",
		cmd.ENV_KEY_FOR_K8S:           "ON",
		cmd.ENV_KEY_FOR_FIREFOX:       "ON",
		cmd.ENV_KEY_USER:              "certowner",
		cmd.ENV_KEY_UID:               "1000",
		cmd.ENV_KEY_GID:               "1000",
		cmd.ENV_KEY_OUTPUT_DIR:        "env-certs",
		cmd.ENV_KEY_EXPIRE_DAYS:       "30",
		cmd.ENV_KEY_CUSTOM_FILE_NAME:  "environment-certificate",
	}
	fromEnvironment := cmd.AppFlags{
		Country:            "US",
		State:              "Environment State",
		Locality:           "ENVIRONMENT LOCALITY",
		Organization:       "Environment Organization",
		OrganizationalUnit: "Environment Unit",
		CommonName:         "Environment Common Name",
		Domains:            "env.example.com,*.env.example.com",
		ForK8s:             "true",
		ForFirefox:         "true",
		User:               "certowner",
		UID:                "1000",
		GID:                "1000",
		OutputDir:          "env-certs",
		ExpireDays:         "30",
		CustomFileName:     "environment-certificate",
	}
	explicitDefaults := []string{
		"--CERT_C=" + define.DEFAULT_COUNTRY,
		"--CERT_ST=" + define.DEFAULT_STATE,
		"--CERT_L=" + define.DEFAULT_LOCALITY,
		"--CERT_O=" + define.DEFAULT_ORGANIZATION,
		"--CERT_OU=" + define.DEFAULT_ORGANIZATIONAL_UNIT,
		"--CERT_CN=" + define.DEFAULT_COMMON_NAME,
		"--CERT_DNS=" + define.DEFAULT_DOMAINS,
		"--FOR_K8S=OFF",
		"--FOR_FIREFOX=OFF",
		"--DIR=" + define.DEFAULT_DIR,
		"--EXPIRE_DAYS=" + define.DEFAULT_EXPIRE_DAYS,
	}
	defaultsWithEnvironmentFilename := defaults
	defaultsWithEnvironmentFilename.CustomFileName = fromEnvironment.CustomFileName
	cliOverrides := fromEnvironment
	cliOverrides.Country = "GB"
	cliOverrides.User = "cliowner"
	cliOverrides.UID = "2000"
	cliOverrides.GID = "2000"
	cliOverrides.OutputDir = "cli-certs"
	cliOverrides.ExpireDays = "90"
	cliOverrides.CustomFileName = "cli-certificate"

	tests := []struct {
		name string
		env  map[string]string
		args []string
		want cmd.AppFlags
	}{
		{"defaults without CLI or environment", nil, nil, defaults},
		{"unspecified CLI preserves environment", environment, nil, fromEnvironment},
		{"explicit defaults override environment", environment, explicitDefaults, defaultsWithEnvironmentFilename},
		{"explicit CLI overrides environment", environment, []string{
			"--CERT_C=gb", "--USER=cliowner", "--UID=2000", "--GID=2000",
			"--DIR=cli-certs", "--EXPIRE_DAYS=90", "--CUSTOM_FILE_NAME=cli-certificate",
		}, cliOverrides},
		{"explicit empty CLI preserves environment", environment, []string{
			"--CERT_C=", "--FOR_K8S=", "--FOR_FIREFOX=", "--DIR=", "--EXPIRE_DAYS=", "--CUSTOM_FILE_NAME=",
		}, fromEnvironment},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for key := range environment {
				t.Setenv(key, "")
			}
			t.Setenv(cmd.ENV_KEY_USER, "certowner")
			t.Setenv(cmd.ENV_KEY_UID, "1000")
			t.Setenv(cmd.ENV_KEY_GID, "1000")
			for key, value := range test.env {
				t.Setenv(key, value)
			}
			flag.CommandLine = flag.NewFlagSet("certs-maker", flag.ContinueOnError)
			os.Args = append([]string{"certs-maker"}, test.args...)
			cmd.ApplyFlags()

			got := cmd.AppFlags{
				Country:            define.CERT_COUNTRY,
				State:              define.CERT_STATE,
				Locality:           define.CERT_LOCALITY,
				Organization:       define.CERT_ORGANIZATION,
				OrganizationalUnit: define.CERT_ORGANIZATIONAL_UNIT,
				CommonName:         define.CERT_COMMON_NAME,
				Domains:            strings.Join(define.CERT_DOMAINS, ","),
				ForK8s:             strconv.FormatBool(define.APP_FOR_K8S),
				ForFirefox:         strconv.FormatBool(define.APP_FOR_FIREFOX),
				User:               define.APP_USER,
				UID:                define.APP_UID,
				GID:                define.APP_GID,
				OutputDir:          define.APP_OUTPUT_DIR,
				ExpireDays:         define.CERT_EXPIRE_DAYS,
				CustomFileName:     define.CUSTOM_FILE_NAME,
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("ApplyFlags() = %+v, want %+v", got, test.want)
			}
			if _, err := os.Stat(define.APP_OUTPUT_DIR); err != nil {
				t.Fatalf("output directory was not created: %v", err)
			}
		})
	}
}

func TestParseFlagsKeepsHelpDefaults(t *testing.T) {
	originalFlags, originalArgs := flag.CommandLine, os.Args
	t.Cleanup(func() {
		flag.CommandLine, os.Args = originalFlags, originalArgs
	})
	flag.CommandLine = flag.NewFlagSet("certs-maker", flag.ContinueOnError)
	os.Args = []string{"certs-maker"}
	parsed := cmd.ParseFlags()
	if parsed != (cmd.AppFlags{}) {
		t.Fatalf("unspecified CLI values = %+v, want empty values for environment fallback", parsed)
	}
	for name, want := range map[string]string{
		cmd.ENV_KEY_COUNTRY:    define.DEFAULT_COUNTRY,
		cmd.ENV_KEY_FOR_K8S:    define.DEFAULT_FOR_K8S,
		cmd.ENV_KEY_OUTPUT_DIR: define.DEFAULT_DIR,
		cmd.ENV_KEY_GID:        define.DEFAULT_GID,
	} {
		if got := flag.Lookup(name).DefValue; got != want {
			t.Errorf("%s help default = %q, want %q", name, got, want)
		}
	}
}
