package spec

type Spec struct {
	Path     string
	Title    string
	Sections map[string]string
}

type Graph struct {
	Nodes map[string]*Node
	Edges []Edge
}

type Node struct {
	ID   string
	Path string
	Type string
}

type Edge struct {
	From string
	To   string
}
