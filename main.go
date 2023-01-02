package main

import (
	"fmt"

	"github.com/soulteary/certs-maker/internal/cmd"
	"github.com/soulteary/certs-maker/internal/generator"
	"github.com/soulteary/certs-maker/internal/version"
)

func main() {
	fmt.Printf("running soulteary/certs-maker %s\n", version.Version)
	cmd.ApplyFlags()
	// TODO print all flags
	generator.Generate()
}
