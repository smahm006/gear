package inventory

type Group struct {
	Name        string
	ParentGroup string
	SubGroup    map[string]*Group
	Hosts       map[string]*Host
	Variables   map[string]interface{}
	Environment map[string]interface{}
}

func NewGroup(name string) *Group {
	return &Group{
		Name:        name,
		Variables:   make(map[string]interface{}),
		Environment: make(map[string]interface{}),
	}
}
