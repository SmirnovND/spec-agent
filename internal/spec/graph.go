package spec

import (
	"path/filepath"
)

func CollectAllReferences(specFiles []string) map[string]bool {
	referenced := map[string]bool{}

	for _, f := range specFiles {
		_, edges, err := ParseDependencies(f)
		if err != nil {
			continue
		}

		for _, e := range edges {
			referenced[e.To] = true
		}
	}

	return referenced
}

func FindRootSpecs(specs []string, referenced map[string]bool) []string {
	var roots []string

	for _, s := range specs {
		abs, _ := filepath.Abs(s)
		if !referenced[abs] {
			roots = append(roots, abs)
		}
	}

	return roots
}

func BuildGraphFromRoots(rootSpecs []string) (*Graph, error) {
	graph := &Graph{
		Nodes: map[string]*Node{},
		Edges: []Edge{},
	}

	visited := map[string]bool{}

	for _, root := range rootSpecs {
		if err := walkSpec(root, graph, visited); err != nil {
			return nil, err
		}
	}

	return graph, nil
}

func walkSpec(specPath string, graph *Graph, visited map[string]bool) error {
	if visited[specPath] {
		return nil
	}
	visited[specPath] = true

	if _, exists := graph.Nodes[specPath]; !exists {
		graph.Nodes[specPath] = &Node{
			ID:   specPath,
			Path: specPath,
			Type: "spec",
		}
	}

	_, edges, err := ParseDependencies(specPath)
	if err != nil {
		return nil
	}

	for _, edge := range edges {
		graph.Edges = append(graph.Edges, Edge{
			From: specPath,
			To:   edge.To,
		})

		if _, exists := graph.Nodes[edge.To]; !exists {
			graph.Nodes[edge.To] = &Node{
				ID:   edge.To,
				Path: edge.To,
				Type: "spec",
			}
		}

		if err := walkSpec(edge.To, graph, visited); err != nil {
			return err
		}
	}

	return nil
}
