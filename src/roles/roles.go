package roles

type Role struct {
	Name      string                 `yaml:"name"`
	Variables map[string]interface{} `yaml:"vars"`
	Tasks     []string               `yaml:"tasks"`
	Handlers  []string               `yaml:"handlers"`
}

func NewRole() *Role {
	return &Role{}
}
