package fn_test

import (
	"testing"

	"github.com/soulteary/certs-maker/internal/fn"
)

func TestIsNotEmptyAndNotDefaultString(t *testing.T) {
	ret := fn.IsNotEmptyAndNotDefaultString("", "d")
	if ret {
		t.Fatal("IsNotEmptyAndNotDefaultString failed")
	}

	ret = fn.IsNotEmptyAndNotDefaultString("d", "d")
	if ret {
		t.Fatal("IsNotEmptyAndNotDefaultString failed")
	}
}

func TestIsBoolString(t *testing.T) {
	ret := fn.IsBoolString("true")
	if !ret {
		t.Fatal("IsBoolString failed")
	}
	ret = fn.IsBoolString("on")
	if !ret {
		t.Fatal("IsBoolString failed")
	}
	ret = fn.IsBoolString("1")
	if !ret {
		t.Fatal("IsBoolString failed")
	}
	ret = fn.IsBoolString("ON")
	if !ret {
		t.Fatal("IsBoolString failed")
	}

	ret = fn.IsBoolString("not-vaild")
	if ret {
		t.Fatal("IsBoolString failed")
	}
}

func TestIsDomainListStringMatch(t *testing.T) {
	ret := fn.IsDomainListStringMatch([]string{"a.com", "b.com"}, "a.com,b.com")
	if !ret {
		t.Fatal("Test IsDomainListStringMatch failed")
	}

	ret = fn.IsDomainListStringMatch([]string{"a.com", "b.com"}, "a.com,b.com,c.com")
	if ret {
		t.Fatal("Test IsDomainListStringMatch failed")
	}
}

func TestIsVaildCountry(t *testing.T) {
	ret := fn.IsVaildCountry("AA")
	if !ret {
		t.Fatal("Test IsVaildCountry failed")
	}

	ret = fn.IsVaildCountry("AA!")
	if ret {
		t.Fatal("Test IsVaildCountry failed")
	}
	ret = fn.IsVaildCountry("")
	if ret {
		t.Fatal("Test IsVaildCountry failed")
	}

	ret = fn.IsVaildCountry("a")
	if ret {
		t.Fatal("Test IsVaildCountry failed")
	}

	ret = fn.IsVaildCountry("a!")
	if ret {
		t.Fatal("Test IsVaildCountry failed")
	}
}

func TestIsStrInArray(t *testing.T) {
	ret := fn.IsStrInArray([]string{"a", "b"}, "a")
	if !ret {
		t.Fatal("Checking string array contains data failed")
	}

	ret = fn.IsStrInArray([]string{"a", "b"}, "a!!")
	if ret {
		t.Fatal("Checking string array contains data failed")
	}
}
