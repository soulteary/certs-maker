/**
 * Copyright (c) 2021-2025 Su Yang (soulteary)
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package cmd_test

import (
	"os"
	"testing"

	"github.com/soulteary/certs-maker/internal/cmd"
	"github.com/soulteary/certs-maker/internal/define"
	"github.com/soulteary/certs-maker/internal/fn"
)

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
	os.Setenv(cmd.ENV_KEY_COUNTRY, "")
}

func TestUpdateDomainOption(t *testing.T) {
	// env: empty, args: "", default: define.DEFAULT_DOMAINS
	ret := cmd.UpdateDomainOption(cmd.ENV_KEY_DOMAINS, "", define.DEFAULT_DOMAINS)
	if !fn.IsDomainListStringMatch(ret, define.DEFAULT_DOMAINS) {
		t.Fatal("UpdateDomainOption failed")
	}

	// env: empty, args: "a.com,b.com", default: define.DEFAULT_DOMAINS
	ret = cmd.UpdateDomainOption(cmd.ENV_KEY_DOMAINS, "a.com,b.com", define.DEFAULT_DOMAINS)
	if fn.IsDomainListStringMatch(ret, define.DEFAULT_DOMAINS) {
		t.Fatal("UpdateDomainOption failed")
	}
	if !fn.IsStrInArray(ret, "a.com") || !fn.IsStrInArray(ret, "b.com") {
		t.Fatal("UpdateDomainOption failed")
	}

	// env: "c.com", args: "", default: define.DEFAULT_DOMAINS
	os.Setenv(cmd.ENV_KEY_DOMAINS, "c.com")
	ret = cmd.UpdateDomainOption(cmd.ENV_KEY_DOMAINS, "", define.DEFAULT_DOMAINS)
	if fn.IsDomainListStringMatch(ret, define.DEFAULT_DOMAINS) {
		t.Fatal("UpdateDomainOption failed")
	}
	if !fn.IsStrInArray(ret, "c.com") {
		t.Fatal("UpdateDomainOption failed")
	}

	// env: "c.com", args: "a.com", default: define.DEFAULT_DOMAINS
	os.Setenv(cmd.ENV_KEY_DOMAINS, "c.com")
	ret = cmd.UpdateDomainOption(cmd.ENV_KEY_DOMAINS, "a.com", define.DEFAULT_DOMAINS)
	if fn.IsDomainListStringMatch(ret, define.DEFAULT_DOMAINS) {
		t.Fatal("UpdateDomainOption failed")
	}
	if !fn.IsStrInArray(ret, "a.com") {
		t.Fatal("UpdateDomainOption failed")
	}

	os.Setenv(cmd.ENV_KEY_DOMAINS, "")
}

func TestSantizeDirPath(t *testing.T) {
	// env: empty, args: "", default: define.DEFAULT_DIR
	ret := cmd.SantizeDirPath(cmd.ENV_KEY_OUTPUT_DIR, "", define.DEFAULT_DIR)
	if ret != define.DEFAULT_DIR {
		t.Fatal("test SantizeDirPath failed")
	}

	// env: empty, args: "/aa", default: define.DEFAULT_DIR
	ret = cmd.SantizeDirPath(cmd.ENV_KEY_OUTPUT_DIR, "/aa", define.DEFAULT_DIR)
	if ret != "aa" {
		t.Fatal("test SantizeDirPath failed")
	}

	// env: empty, args: "./aaa", default: define.DEFAULT_DIR
	ret = cmd.SantizeDirPath(cmd.ENV_KEY_OUTPUT_DIR, "./aaa", define.DEFAULT_DIR)
	if ret != "aaa" {
		t.Fatal("test SantizeDirPath failed")
	}

	// env: empty, args: "../aaa", default: define.DEFAULT_DIR
	ret = cmd.SantizeDirPath(cmd.ENV_KEY_OUTPUT_DIR, "../aaa", define.DEFAULT_DIR)
	if ret != "aaa" {
		t.Fatal("test SantizeDirPath failed")
	}

	// env: empty, args: "....//abc", default: define.DEFAULT_DIR
	ret = cmd.SantizeDirPath(cmd.ENV_KEY_OUTPUT_DIR, "....//abc", define.DEFAULT_DIR)
	if ret != "abc" {
		t.Fatal("test SantizeDirPath failed")
	}

	// env: empty, args: "abcd/abc", default: define.DEFAULT_DIR
	ret = cmd.SantizeDirPath(cmd.ENV_KEY_OUTPUT_DIR, "abcd/abc", define.DEFAULT_DIR)
	if ret != "abcd/abc" {
		t.Fatal("test SantizeDirPath failed")
	}

	// env: empty, args: "././././abcd", default: define.DEFAULT_DIR
	ret = cmd.SantizeDirPath(cmd.ENV_KEY_OUTPUT_DIR, "././././abcd", define.DEFAULT_DIR)
	if ret != "abcd" {
		t.Fatal("test SantizeDirPath failed")
	}

	// env: empty, args: "...../aaa/a././aa", default: define.DEFAULT_DIR
	ret = cmd.SantizeDirPath(cmd.ENV_KEY_OUTPUT_DIR, "...../aaa/a././aa", define.DEFAULT_DIR)
	if ret != "aaa/aaa" {
		t.Fatal("test SantizeDirPath failed")
	}
}

func TestUpdateBoolOption(t *testing.T) {

	// env: empty, args: false, default: false
	ret := cmd.UpdateBoolOption("TEST_KEY", "false", "false")
	if ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: empty, args: false, default: true
	ret = cmd.UpdateBoolOption("TEST_KEY", "false", "true")
	if ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: empty, args: true, default: true
	ret = cmd.UpdateBoolOption("TEST_KEY", "true", "true")
	if !ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: empty, args: false, default: true
	ret = cmd.UpdateBoolOption("TEST_KEY", "false", "true")
	if ret {
		t.Fatal("UpdateBoolOption failed")
	}

	// env: on, args: false, default: false
	os.Setenv("TEST_KEY", "on")
	ret = cmd.UpdateBoolOption("TEST_KEY", "false", "false")
	if !ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: on, args: true, default: false
	os.Setenv("TEST_KEY", "on")
	ret = cmd.UpdateBoolOption("TEST_KEY", "true", "false")
	if !ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: on, args: true, default: true
	os.Setenv("TEST_KEY", "on")
	ret = cmd.UpdateBoolOption("TEST_KEY", "true", "true")
	if !ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: on, args: true, default: false
	os.Setenv("TEST_KEY", "on")
	ret = cmd.UpdateBoolOption("TEST_KEY", "true", "false")
	if !ret {
		t.Fatal("UpdateBoolOption failed")
	}

	// env: off, args: false, default: false
	os.Setenv("TEST_KEY", "off")
	ret = cmd.UpdateBoolOption("TEST_KEY", "false", "false")
	if ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: off, args: true, default: false
	os.Setenv("TEST_KEY", "on")
	ret = cmd.UpdateBoolOption("TEST_KEY", "true", "false")
	if !ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: off, args: true, default: true
	os.Setenv("TEST_KEY", "on")
	ret = cmd.UpdateBoolOption("TEST_KEY", "true", "true")
	if !ret {
		t.Fatal("UpdateBoolOption failed")
	}
	// env: on, args: true, default: false
	os.Setenv("TEST_KEY", "on")
	ret = cmd.UpdateBoolOption("TEST_KEY", "true", "false")
	if !ret {
		t.Fatal("UpdateBoolOption failed")
	}

	os.Setenv("TEST_KEY", "")
}
