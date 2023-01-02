package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/soulteary/certs-maker/internal/cmd"
	"github.com/soulteary/certs-maker/internal/fn"
	"github.com/soulteary/certs-maker/internal/generator"
	"github.com/soulteary/certs-maker/internal/version"
)

func init() {
	prepare := filepath.Join(".", "ssl")
	os.MkdirAll(prepare, os.ModePerm)
}

func main() {
	fmt.Printf("running soulteary/certs-maker %s\n", version.Version)
	cmd.ApplyFlags()

	shell := generator.GenerateConfFile()
	fn.Execute(shell)
	generator.TryAdjustPermissions()
}
