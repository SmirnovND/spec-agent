package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindSpecsNearRoots_ExcludeDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	mustWriteSpec(t, filepath.Join(tmpDir, "cmd", "public", "user_controller.md"))
	mustWriteSpec(t, filepath.Join(tmpDir, "cmd", "staticlint", "lint.md"))
	mustWriteSpec(t, filepath.Join(tmpDir, "cmd", "server", "server.md"))

	roots := []string{filepath.Join(tmpDir, "cmd")}
	exclude := []string{
		filepath.Join(tmpDir, "cmd", "staticlint"),
		filepath.Join(tmpDir, "cmd", "server"),
	}

	specs, err := findSpecsNearRoots(roots, exclude)
	if err != nil {
		t.Fatalf("findSpecsNearRoots вернул ошибку: %v", err)
	}

	if len(specs) != 1 {
		t.Fatalf("ожидался 1 spec после exclude, получено %d: %#v", len(specs), specs)
	}

	got := filepath.Base(specs[0])
	if got != "user_controller.md" {
		t.Fatalf("ожидался user_controller.md, получен %s", got)
	}
}

func mustWriteSpec(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("не удалось создать директорию: %v", err)
	}

	content := `<!-- SPEC:FILE=true -->
<!-- SPEC:ID=test/spec -->
<!-- SPEC:KIND=controller -->
<!-- SPEC:MENU=true -->
# TestSpec
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("не удалось записать spec: %v", err)
	}
}
