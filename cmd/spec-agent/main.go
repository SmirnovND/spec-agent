package main

import (
	"os"
	"runtime/debug"
	"strings"

	"github.com/SmirnovND/spec-agent/v2/internal/cli"
)

var version = "dev"
var readBuildInfo = debug.ReadBuildInfo

func main() {
	cli.Version = detectVersion()

	// Support common CLI convention: `spec-agent --version` and `spec-agent -v`.
	if len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		os.Args = []string{os.Args[0], "version"}
	}

	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}

func detectVersion() string {
	v := strings.TrimSpace(version)
	if v != "" && v != "dev" {
		return v
	}

	if info, ok := readBuildInfo(); ok {
		mv := strings.TrimSpace(info.Main.Version)
		if mv != "" && mv != "(devel)" {
			return mv
		}
	}

	return "dev"
}
