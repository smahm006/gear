package inventory

type Group struct {
	Name         string
	ParentGroups map[string]*Group
	SubGroups    map[string]*Group
	Hosts        map[string]*Host
	Variables    map[string]string
	Environment  map[string]string
}

func NewGroup(name string) *Group {
	return &Group{
		Name:         name,
		ParentGroups: make(map[string]*Group),
		SubGroups:    make(map[string]*Group),
		Hosts:        make(map[string]*Host),
		Variables:    make(map[string]string),
		Environment:  make(map[string]string),
	}
}

func (g *Group) GetEnv(name string) string {
	env := g.Environment[name]
	return env
}

func (g *Group) GetHost(name string) (*Host, bool) {
	host, ok := g.Hosts[name]
	return host, ok
}
