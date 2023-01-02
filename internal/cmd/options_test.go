package cmd_test

import (
	"os"
	"testing"

	"github.com/soulteary/certs-maker/internal/cmd"
	"github.com/soulteary/certs-maker/internal/define"
)

func TestUpdateBoolOption(t *testing.T) {

	// env: empty, args: false, default: false
	ret := cmd.UpdateBoolOption("TEST_KEY", false, false)
	if ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: empty, args: false, default: true
	ret = cmd.UpdateBoolOption("TEST_KEY", false, true)
	if ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: empty, args: true, default: true
	ret = cmd.UpdateBoolOption("TEST_KEY", true, true)
	if !ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: empty, args: false, default: true
	ret = cmd.UpdateBoolOption("TEST_KEY", false, true)
	if ret {
		t.Fatal("UpdateBoolOption failed")
	}

	// env: on, args: false, default: false
	os.Setenv("TEST_KEY", "on")
	ret = cmd.UpdateBoolOption("TEST_KEY", false, false)
	if !ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: on, args: true, default: false
	os.Setenv("TEST_KEY", "on")
	ret = cmd.UpdateBoolOption("TEST_KEY", true, false)
	if !ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: on, args: true, default: true
	os.Setenv("TEST_KEY", "on")
	ret = cmd.UpdateBoolOption("TEST_KEY", true, true)
	if !ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: on, args: true, default: false
	os.Setenv("TEST_KEY", "on")
	ret = cmd.UpdateBoolOption("TEST_KEY", true, false)
	if !ret {
		t.Fatal("UpdateBoolOption failed")
	}

	// env: off, args: false, default: false
	os.Setenv("TEST_KEY", "off")
	ret = cmd.UpdateBoolOption("TEST_KEY", false, false)
	if ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: off, args: true, default: false
	os.Setenv("TEST_KEY", "on")
	ret = cmd.UpdateBoolOption("TEST_KEY", true, false)
	if !ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: off, args: true, default: true
	os.Setenv("TEST_KEY", "on")
	ret = cmd.UpdateBoolOption("TEST_KEY", true, true)
	if !ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: on, args: true, default: false
	os.Setenv("TEST_KEY", "on")
	ret = cmd.UpdateBoolOption("TEST_KEY", true, false)
	if !ret {
		t.Fatal("UpdateBoolOption failed")
	}

	os.Setenv("TEST_KEY", "")
}

func TestUpdateStringOption(t *testing.T) {
	// env: empty, args:"a", default:"d"
	ret := cmd.UpdateStringOption("TEST_KEY", "a", "d")
	if ret != "a" {
		t.Fatal("UpdateStringOption failed")
	}
	// env: empty, args:"", default:"d"
	ret = cmd.UpdateStringOption("TEST_KEY", "", "d")
	if ret != "d" {
		t.Fatal("UpdateStringOption failed")
	}
	// env: empty, args:"a", default:""
	ret = cmd.UpdateStringOption("TEST_KEY", "a", "")
	if ret != "a" {
		t.Fatal("UpdateStringOption failed")
	}
	// env: empty, args:"", default:""
	ret = cmd.UpdateStringOption("TEST_KEY", "", "")
	if ret != "" {
		t.Fatal("UpdateStringOption failed")
	}

	os.Setenv("TEST_KEY", "e")
	// env: "e", args:"a", default:"d"
	ret = cmd.UpdateStringOption("TEST_KEY", "a", "d")
	if ret != "a" {
		t.Fatal("UpdateStringOption failed")
	}
	// env: "e", args:"", default:"d"
	ret = cmd.UpdateStringOption("TEST_KEY", "", "d")
	if ret != "e" {
		t.Fatal("UpdateStringOption failed")
	}
	// env: "e", args:"a", default:""
	ret = cmd.UpdateStringOption("TEST_KEY", "a", "")
	if ret != "a" {
		t.Fatal("UpdateStringOption failed")
	}
	// env: "e", args:"", default:""
	ret = cmd.UpdateStringOption("TEST_KEY", "", "")
	if ret != "e" {
		t.Fatal("UpdateStringOption failed")
	}
	os.Setenv("TEST_KEY", "")
}

func TestUpdateCountryOption(t *testing.T) {
	// env: empty, args: "", default: define.DEFAULT_COUNTRY
	ret := cmd.UpdateCountryOption(cmd.ENV_KEY_COUNTRY, "", define.DEFAULT_COUNTRY)
	if ret != define.DEFAULT_COUNTRY {
		t.Fatal("UpdateCountryOption failed")
	}

	// env: empty, args: "xx", default: define.DEFAULT_COUNTRY
	ret = cmd.UpdateCountryOption(cmd.ENV_KEY_COUNTRY, "xx", define.DEFAULT_COUNTRY)
	if ret != "XX" {
		t.Fatal("UpdateCountryOption failed")
	}

	// env: empty, args: "!x!x!", default: define.DEFAULT_COUNTRY
	ret = cmd.UpdateCountryOption(cmd.ENV_KEY_COUNTRY, "!x!x!", define.DEFAULT_COUNTRY)
	if ret != define.DEFAULT_COUNTRY {
		t.Fatal("UpdateCountryOption failed")
	}

	// env: empty, args: "x", default: define.DEFAULT_COUNTRY
	ret = cmd.UpdateCountryOption(cmd.ENV_KEY_COUNTRY, "x", define.DEFAULT_COUNTRY)
	if ret != define.DEFAULT_COUNTRY {
		t.Fatal("UpdateCountryOption failed")
	}

	// env: "xx", args: empty, default: define.DEFAULT_COUNTRY
	os.Setenv(cmd.ENV_KEY_COUNTRY, "xx")
	ret = cmd.UpdateCountryOption(cmd.ENV_KEY_COUNTRY, "", define.DEFAULT_COUNTRY)
	if ret != "XX" {
		t.Fatal("UpdateCountryOption failed")
	}
}
