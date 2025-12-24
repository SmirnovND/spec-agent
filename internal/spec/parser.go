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

	content := string(data)
	lines := strings.Split(content, "\n")

	spec := &Spec{
		Path:     path,
		Content:  content,
		Sections: map[string]string{},
		Links:    []SpecLink{},
	}

	var current string
	re := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+\.md)\)`)

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
		
		matches := re.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			spec.Links = append(spec.Links, SpecLink{
				Title: match[1],
				Path:  match[2],
			})
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
