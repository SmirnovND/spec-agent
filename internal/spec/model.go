package spec

type Spec struct {
	Path     string
	Title    string
	Sections map[string]string
	Content  string
	Links    []SpecLink
	Meta     SpecMeta
}

type SpecLink struct {
	Title  string
	Path   string
	Anchor string
}

type SpecMeta struct {
	IsSpecFile bool
	ID         string
	Kind       string
	Menu       bool
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

type ExportTree struct {
	Specs     map[string]*Spec
	RootSpecs []string
}
