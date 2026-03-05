package main

import (
	"os"

	"github.com/SmirnovND/spec-agent/v3/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
