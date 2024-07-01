package inventory

type Group struct {
	Name        string
	ParentGroup map[string]*Group
	SubGroup    map[string]*Group
	Hosts       map[string]*Host
	Variables   map[string]interface{}
	Environment map[string]interface{}
}

func NewGroup(name string) *Group {
	return &Group{
		Name:        name,
		ParentGroup: make(map[string]*Group),
		SubGroup:    make(map[string]*Group),
		Variables:   make(map[string]interface{}),
		Environment: make(map[string]interface{}),
	}
}
