package main

import (
	"os"

	"github.com/SmirnovND/spec-agent/v2/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
