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
