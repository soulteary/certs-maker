/**
 * Copyright (c) 2021-2025 Su Yang (soulteary)
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package fn_test

import (
	"testing"

	"github.com/soulteary/certs-maker/internal/fn"
)

func TestExecute(t *testing.T) {
	out, err := fn.Execute("ls")
	if err != nil {
		t.Fatal("test Execute failed")
	}
	if len(out) == 0 {
		t.Fatal("test Execute failed")
	}

	out, err = fn.Execute("not-exist-command")
	if err == nil {
		t.Fatal("test Execute failed")
	}
	if len(out) == 0 {
		t.Fatal("test Execute failed")
	}
}
