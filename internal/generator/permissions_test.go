package generator_test

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"testing"

	"github.com/soulteary/certs-maker/internal/define"
	"github.com/soulteary/certs-maker/internal/generator"
)

func TestGeneratePermissionsFixsCmds(t *testing.T) {
	define.APP_USER = "testuser"
	define.APP_UID = "12345"
	define.APP_GID = "09876"
	define.APP_OUTPUT_DIR = "./ssl"

	ret := generator.GeneratePermissionsFixsCmds()

	check1 := strings.Contains(ret, fmt.Sprintf("addgroup -g %s %s", define.APP_GID, define.APP_USER))
	if !check1 {
		t.Fatal("Test GeneratePermissionsFixsCmds failed")
	}

	check2 := strings.Contains(ret, fmt.Sprintf(`adduser -g "" -G %s -H -D -u %s %s`, define.APP_USER, define.APP_UID, define.APP_USER))
	if !check2 {
		t.Fatal("Test GeneratePermissionsFixsCmds failed")
	}

	check3 := strings.Contains(ret, fmt.Sprintf(`chown -R %s:%s %s`, define.APP_USER, define.APP_USER, define.APP_OUTPUT_DIR))
	if !check3 {
		t.Fatal("Test GeneratePermissionsFixsCmds failed")
	}

	check4 := strings.Contains(ret, fmt.Sprintf(`chmod -R a+r %s`, define.APP_OUTPUT_DIR))
	if !check4 {
		t.Fatal("Test GeneratePermissionsFixsCmds failed")
	}

	define.APP_USER = ""
	define.APP_UID = ""
	define.APP_GID = ""
	define.APP_OUTPUT_DIR = ""
	ret = generator.GeneratePermissionsFixsCmds()
	if ret != "" {
		t.Fatal("Test GeneratePermissionsFixsCmds failed")
	}
}

func TestTryToFixPermissions(t *testing.T) {
	define.APP_USER = ""
	define.APP_UID = ""
	define.APP_GID = ""
	define.APP_OUTPUT_DIR = ""

	err := generator.TryToFixPermissions()
	if err != nil {
		t.Fatal("Test TryToFixPermissions failed")
	}

	os.MkdirAll(define.DEFAULT_DIR, os.ModePerm)

	user, err := user.Current()
	fmt.Println(user.Username, user.Gid, user.Uid)
	if err != nil {
		t.Fatal("Test TryToFixPermissions failed, get current user failed")
	}

	define.APP_USER = user.Username
	define.APP_UID = user.Uid
	define.APP_GID = user.Gid
	define.APP_OUTPUT_DIR = define.DEFAULT_DIR

	err = generator.TryToFixPermissions()
	if err != nil {
		fmt.Println(err)
		t.Fatal("Test TryToFixPermissions failed")
	}
	os.RemoveAll(define.DEFAULT_DIR)
}
