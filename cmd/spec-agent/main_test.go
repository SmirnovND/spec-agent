package main

import (
	"runtime/debug"
	"testing"
)

func TestDetectVersion_UsesInjectedVersion(t *testing.T) {
	prevVersion := version
	prevRead := readBuildInfo
	t.Cleanup(func() {
		version = prevVersion
		readBuildInfo = prevRead
	})

	version = "v2.0.4"
	readBuildInfo = func() (info *debug.BuildInfo, ok bool) {
		return nil, false
	}

	if got := detectVersion(); got != "v2.0.4" {
		t.Fatalf("detectVersion() = %q, want %q", got, "v2.0.4")
	}
}

func TestDetectVersion_FallsBackToBuildInfo(t *testing.T) {
	prevVersion := version
	prevRead := readBuildInfo
	t.Cleanup(func() {
		version = prevVersion
		readBuildInfo = prevRead
	})

	version = "dev"
	readBuildInfo = func() (info *debug.BuildInfo, ok bool) {
		return &debug.BuildInfo{
			Main: debug.Module{
				Version: "v2.0.3",
			},
		}, true
	}

	if got := detectVersion(); got != "v2.0.3" {
		t.Fatalf("detectVersion() = %q, want %q", got, "v2.0.3")
	}
}

