package spec

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var specMetaLineRE = regexp.MustCompile(`<!--\s*SPEC:([A-Z]+)=([^\s]+)\s*-->`)

func ParseMeta(content string) SpecMeta {
	meta := SpecMeta{}
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		matches := specMetaLineRE.FindStringSubmatch(line)
		if len(matches) != 3 {
			continue
		}

		key := matches[1]
		value := matches[2]
		switch key {
		case "FILE":
			meta.IsSpecFile = strings.EqualFold(value, "true")
		case "ID":
			meta.ID = value
		case "KIND":
			meta.Kind = value
		case "MENU":
			meta.Menu = strings.EqualFold(value, "true")
		}
	}

	return meta
}

func ParseMetaFromFile(path string) (SpecMeta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return SpecMeta{}, err
	}
	return ParseMeta(string(data)), nil
}

func IsSpecFile(path string) (bool, error) {
	meta, err := ParseMetaFromFile(path)
	if err != nil {
		return false, err
	}
	return meta.IsSpecFile, nil
}

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
		Meta:     ParseMeta(content),
	}

	var current string
	re := regexp.MustCompile(`\[([^\]]+)\]\(([^)#]+\.md)(#[^)]+)?\)`)

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
			anchor := ""
			if len(match) > 3 {
				anchor = strings.TrimPrefix(match[3], "#")
			}
			spec.Links = append(spec.Links, SpecLink{
				Title:  match[1],
				Path:   match[2],
				Anchor: anchor,
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
	re := regexp.MustCompile(`\[.*?\]\(([^)#]+\.md)(#[^)]+)?\)`)

	dir := filepath.Dir(specPath)

	for _, section := range spec.Sections {
		matches := re.FindAllStringSubmatch(section, -1)
		for _, match := range matches {
			refPath := match[1]
			absPath := filepath.Join(dir, refPath)
			normalized, _ := filepath.Abs(absPath)
			ok, err := IsSpecFile(normalized)
			if err != nil || !ok {
				continue
			}

			edges = append(edges, Edge{
				From: specPath,
				To:   normalized,
			})
		}
	}

	return spec, edges, nil
}
