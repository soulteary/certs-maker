package generator

import (
	"fmt"

	"github.com/soulteary/certs-maker/internal/define"
	"github.com/soulteary/certs-maker/internal/fn"
)

func TryToFixPermissions() {
	if define.APP_USER == "" ||
		define.APP_UID == "" ||
		define.APP_GID == "" {
		return
	}

	cmdCreateGroup := fmt.Sprintf(`addgroup -g %s %s`, define.APP_GID, define.APP_USER)
	cmdCreateUser := fmt.Sprintf(`adduser -g "" -G %s -H -D -u %s %s`, define.APP_USER, define.APP_UID, define.APP_USER)
	cmdChangeOwner := fmt.Sprintf(`chown -R %s:%s %s`, define.APP_USER, define.APP_USER, define.APP_OUTPUT_DIR)
	cmdChangeMod := fmt.Sprintf(`chmod -R a+r %s`, define.APP_OUTPUT_DIR)

	fn.Execute(fmt.Sprintf("%s\n%s\n%s\n%s\n", cmdCreateGroup, cmdCreateUser, cmdChangeOwner, cmdChangeMod))
}
