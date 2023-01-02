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
