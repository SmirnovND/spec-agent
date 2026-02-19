package spec

import (
	"strings"
	"testing"
)

func TestGenerateIndexHTML_MenuIncludesOnlyControllerAndCommand(t *testing.T) {
	specs := map[string]*Spec{
		"/tmp/controller.md": {
			Path:  "/tmp/controller.md",
			Title: "UserController",
			Meta: SpecMeta{
				Menu: true,
				Kind: "controller",
			},
		},
		"/tmp/command.md": {
			Path:  "/tmp/command.md",
			Title: "SyncUsersCommand",
			Meta: SpecMeta{
				Menu: true,
				Kind: "command",
			},
		},
		"/tmp/service.md": {
			Path:  "/tmp/service.md",
			Title: "UserService",
			Meta: SpecMeta{
				Menu: true,
				Kind: "service",
			},
		},
		"/tmp/no-menu-controller.md": {
			Path:  "/tmp/no-menu-controller.md",
			Title: "InternalController",
			Meta: SpecMeta{
				Menu: false,
				Kind: "controller",
			},
		},
	}

	html := generateIndexHTML(specs, &Graph{})

	if !strings.Contains(html, "UserController") {
		t.Fatalf("controller должен быть в меню")
	}
	if !strings.Contains(html, "SyncUsersCommand") {
		t.Fatalf("command должен быть в меню")
	}
	if strings.Contains(html, "UserService") {
		t.Fatalf("service не должен быть в меню")
	}
	if strings.Contains(html, "InternalController") {
		t.Fatalf("spec с MENU=false не должен быть в меню")
	}
}
