/**
 * Copyright (c) 2021-2025 Su Yang (soulteary)
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package fn_test

import (
	"reflect"
	"testing"

	"github.com/soulteary/certs-maker/internal/fn"
)

func TestGetUniqDomains(t *testing.T) {
	ret := fn.GetUniqDomains([]string{"a.com", "b.com", "c.com", "a.com"})
	if !reflect.DeepEqual(ret, []string{"a.com", "b.com", "c.com"}) {
		t.Fatal("test GetUniqDomains failed")
	}

	ret = fn.GetUniqDomains([]string{"a.com", "b.com", "c.com"})
	if !reflect.DeepEqual(ret, []string{"a.com", "b.com", "c.com"}) {
		t.Fatal("test GetUniqDomains failed")
	}

	ret = fn.GetUniqDomains([]string{"a.com", "b.com", "c.com", "a.com", "a.com"})
	if !reflect.DeepEqual(ret, []string{"a.com", "b.com", "c.com"}) {
		t.Fatal("test GetUniqDomains failed")
	}

	ret = fn.GetUniqDomains([]string{})
	if !reflect.DeepEqual(ret, []string{}) {
		t.Fatal("test GetUniqDomains failed")
	}
}

func TestGetDomainsByString(t *testing.T) {
	ret := fn.GetDomainsByString("a.com,b.com")
	if !reflect.DeepEqual(ret, []string{"a.com", "b.com"}) {
		t.Fatal("test GetDomainsByString failed")
	}

	ret = fn.GetDomainsByString("a.com,b.com,a.com")
	if !reflect.DeepEqual(ret, []string{"a.com", "b.com"}) {
		t.Fatal("test GetDomainsByString failed")
	}

	ret = fn.GetDomainsByString("A.com,b.com,a.com")
	if !reflect.DeepEqual(ret, []string{"a.com", "b.com"}) {
		t.Fatal("test GetDomainsByString failed")
	}

	ret = fn.GetDomainsByString("A.com,B.com,a.com")
	if !reflect.DeepEqual(ret, []string{"a.com", "b.com"}) {
		t.Fatal("test GetDomainsByString failed")
	}

	ret = fn.GetDomainsByString("192.168.1.1,10.11.12.13")
	if !reflect.DeepEqual(ret, []string{"192.168.1.1", "10.11.12.13"}) {
		t.Fatal("test GetDomainsByString failed")
	}

	ret = fn.GetDomainsByString("2001:db8::1,2001:db8::2")
	if !reflect.DeepEqual(ret, []string{"2001:db8::1", "2001:db8::2"}) {
		t.Fatal("test GetDomainsByString failed")
	}

	ret = fn.GetDomainsByString("2001:db8::1,2001:db8::2,a.com,192.168.1.1,A.COM")
	if !reflect.DeepEqual(ret, []string{"2001:db8::1", "2001:db8::2", "a.com", "192.168.1.1"}) {
		t.Fatal("test GetDomainsByString failed")
	}

	ret = fn.GetDomainsByString("2001:db8::1,2001:db8::2,a.com,192.168.1.1,A.COM,*.a.b.com")
	if !reflect.DeepEqual(ret, []string{"2001:db8::1", "2001:db8::2", "a.com", "192.168.1.1", "*.a.b.com"}) {
		t.Fatal("test GetDomainsByString failed")
	}

	ret = fn.GetDomainsByString("")
	if !reflect.DeepEqual(ret, []string{}) {
		t.Fatal("test GetDomainsByString failed")
	}
}

func TestGetDomainName(t *testing.T) {
	ret := fn.GetDomainName("a.b.com")
	if ret != "a.b.com" {
		t.Fatal("test GetDomainName failed")
	}

	ret = fn.GetDomainName("a.b.c.com")
	if ret != "a.b.c.com" {
		t.Fatal("test GetDomainName failed")
	}

	ret = fn.GetDomainName("com.cc")
	if ret != "com.cc" {
		t.Fatal("test GetDomainName failed")
	}

	ret = fn.GetDomainName("hostname")
	if ret != "hostname" {
		t.Fatal("test GetDomainName failed")
	}

	ret = fn.GetDomainName("")
	if ret != "" {
		t.Fatal("test GetDomainName failed")
	}
}
