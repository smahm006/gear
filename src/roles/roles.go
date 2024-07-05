package roles

type Role struct {
	Name     string   `yaml:"name"`
	Tasks    []string `yaml:"tasks"`
	Handlers []string `yaml:"handlers"`
}

func NewRole() *Role {
	return &Role{}
}
