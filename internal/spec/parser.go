package spec

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ParseFile(path string) (*Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	spec := &Spec{
		Path:     path,
		Sections: map[string]string{},
	}

	var current string

	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			current = strings.TrimPrefix(line, "## ")
			spec.Sections[current] = ""
			continue
		}
		if strings.HasPrefix(line, "# ") && spec.Title == "" {
			spec.Title = strings.TrimPrefix(line, "# ")
		}
		if current != "" {
			spec.Sections[current] += line + "\n"
		}
	}

	return spec, nil
}

func ParseDependencies(specPath string) (*Spec, []Edge, error) {
	spec, err := ParseFile(specPath)
	if err != nil {
		return nil, nil, err
	}

	edges := []Edge{}
	re := regexp.MustCompile(`\[.*?\]\(([^)]+\.md)\)`)

	dir := filepath.Dir(specPath)

	for _, section := range spec.Sections {
		matches := re.FindAllStringSubmatch(section, -1)
		for _, match := range matches {
			refPath := match[1]
			absPath := filepath.Join(dir, refPath)
			normalized, _ := filepath.Abs(absPath)

			edges = append(edges, Edge{
				From: specPath,
				To:   normalized,
			})
		}
	}

	return spec, edges, nil
}
