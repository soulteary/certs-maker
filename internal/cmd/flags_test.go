package cmd_test

import (
	"os"
	"reflect"
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

	if define.APP_FOR_K8S != define.DEFAULT_FOR_K8S {
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
}
