package main

import (
	"os"

	"github.com/SmirnovND/spec-agent/v2/internal/cli"
)

var version = "dev"

func main() {
	cli.Version = version

	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
