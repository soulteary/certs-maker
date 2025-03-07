/**
 * Copyright (c) 2021-2025 Su Yang (soulteary)
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package generator_test

import (
	"errors"
	"os"
	"testing"

	"github.com/soulteary/certs-maker/internal/define"
	"github.com/soulteary/certs-maker/internal/generator"
)

func TestGenerate(t *testing.T) {
	define.APP_FOR_K8S = false
	define.CERT_DOMAINS = append(define.CERT_DOMAINS, "abc.com")
	define.APP_OUTPUT_DIR = "./"

	generator.Generate()

	if _, err := os.Stat("./abc.com.pem.key"); errors.Is(err, os.ErrNotExist) {
		t.Fatal("test MakeCerts failed")
	}
	os.Remove("./abc.com.key")
	if _, err := os.Stat("./abc.com.conf"); errors.Is(err, os.ErrNotExist) {
		t.Fatal("test MakeCerts failed")
	}
	os.Remove("./abc.com.conf")
}
