/**
 * Copyright (c) 2021-2025 Su Yang (soulteary)
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package generator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/soulteary/certs-maker/internal/generator"
)

func TestGenerate(t *testing.T) {
	dir := prepareCertificateTest(t)
	if err := generator.Generate(); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	for _, suffix := range []string{".conf", ".pem.key", ".pem.crt"} {
		if _, err := os.Stat(filepath.Join(dir, "expiry.example.com"+suffix)); err != nil {
			t.Fatalf("expected generated %s: %v", suffix, err)
		}
	}
}
