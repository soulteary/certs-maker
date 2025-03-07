/**
 * Copyright (c) 2021-2025 Su Yang (soulteary)
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"fmt"

	"github.com/soulteary/certs-maker/internal/cmd"
	"github.com/soulteary/certs-maker/internal/generator"
	"github.com/soulteary/certs-maker/internal/version"
)

func main() {
	fmt.Printf("[soulteary/certs-maker] %s\n\n", version.Version)
	cmd.ApplyFlags()
	generator.Generate()
}
