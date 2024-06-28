package main

import (
	"os"
	"testing"

	"github.com/soulteary/certs-maker/internal/cmd"
	"github.com/soulteary/certs-maker/internal/generator"
	"github.com/soulteary/certs-maker/internal/version"
)

func TestMain(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "--CERT_C", "CN", "--CERT_ST", "BJ"}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = oldStdout

	output := make([]byte, 1024)
	n, _ := r.Read(output)
	outputStr := string(output[:n])

	expectedVersionInfo := "[soulteary/certs-maker] " + version.Version
	if outputStr[:len(expectedVersionInfo)] != expectedVersionInfo {
		t.Errorf("Expected output to start with %q, but got %q", expectedVersionInfo, outputStr[:len(expectedVersionInfo)])
	}
}

func TestApplyFlags(t *testing.T) {
	cmd.ApplyFlags()
}

func TestGenerate(t *testing.T) {
	generator.Generate()
}
