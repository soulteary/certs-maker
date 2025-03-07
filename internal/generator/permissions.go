/**
 * Copyright (c) 2021-2025 Su Yang (soulteary)
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package generator

import (
	"fmt"

	"github.com/soulteary/certs-maker/internal/define"
	"github.com/soulteary/certs-maker/internal/fn"
)

// TODO Use Golang to implement functions
func GeneratePermissionsFixsCmds() string {
	if define.APP_USER == "" || define.APP_UID == "" || define.APP_GID == "" {
		return ""
	}
	cmdCreateGroup := fmt.Sprintf(`addgroup -g %s %s`, define.APP_GID, define.APP_USER)
	cmdCreateUser := fmt.Sprintf(`adduser -g "" -G %s -H -D -u %s %s`, define.APP_USER, define.APP_UID, define.APP_USER)
	cmdChangeOwner := fmt.Sprintf(`chown -R %s:%s %s`, define.APP_USER, define.APP_USER, define.APP_OUTPUT_DIR)
	cmdChangeMod := fmt.Sprintf(`chmod -R a+r %s`, define.APP_OUTPUT_DIR)
	return fmt.Sprintf("%s\n%s\n%s\n%s\n", cmdCreateGroup, cmdCreateUser, cmdChangeOwner, cmdChangeMod)
}

func TryToFixPermissions() error {
	shell := GeneratePermissionsFixsCmds()
	if shell == "" {
		return nil
	}
	_, err := fn.Execute(shell)
	return err
}
