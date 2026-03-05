package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestVersionCommand_PrintsInjectedVersion(t *testing.T) {
	prev := Version
	t.Cleanup(func() { Version = prev })

	Version = "v2.0.3"

	out := bytes.NewBuffer(nil)
	rootCmd.SetOut(out)
	rootCmd.SetErr(io.Discard)
	rootCmd.SetArgs([]string{"version"})
	t.Cleanup(func() {
		rootCmd.SetOut(os.Stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs(nil)
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("version command returned error: %v", err)
	}

	if got := strings.TrimSpace(out.String()); got != "v2.0.3" {
		t.Fatalf("unexpected version output: got %q, want %q", got, "v2.0.3")
	}
}

func TestVersionCommand_FallbackToDev(t *testing.T) {
	prev := Version
	t.Cleanup(func() { Version = prev })

	Version = ""

	out := bytes.NewBuffer(nil)
	rootCmd.SetOut(out)
	rootCmd.SetErr(io.Discard)
	rootCmd.SetArgs([]string{"version"})
	t.Cleanup(func() {
		rootCmd.SetOut(os.Stdout)
		rootCmd.SetErr(os.Stderr)
		rootCmd.SetArgs(nil)
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("version command returned error: %v", err)
	}

	if got := strings.TrimSpace(out.String()); got != "dev" {
		t.Fatalf("unexpected version output: got %q, want %q", got, "dev")
	}
}
