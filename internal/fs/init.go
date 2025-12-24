package fs

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed assets
var embeddedAssets embed.FS

func InitSpecAgent() error {
	if err := os.MkdirAll(".spec_agent", 0755); err != nil {
		return err
	}
	if err := os.MkdirAll("spec_changes", 0755); err != nil {
		return err
	}

	config := `roots:
  - internal/controllers
  - cmd
`

	if err := os.WriteFile(".spec_agent/config.yaml", []byte(config), 0644); err != nil {
		return err
	}

	if err := copyAssetsToSpecAgent(); err != nil {
		return err
	}

	return nil
}

func copyAssetsToSpecAgent() error {
	return fs.WalkDir(embeddedAssets, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == "assets" {
			return nil
		}

		relPath := path[len("assets/"):]
		destPath := filepath.Join(".spec_agent", relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		data, err := embeddedAssets.ReadFile(path)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(destPath, data, 0644)
	})
}
