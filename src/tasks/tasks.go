package tasks

type Tasks struct {
	Tag  string
	Name string `yaml:"name"`
	Cmd  string `yaml:"cmd"`
	Save string `yaml:"save"`
	With struct {
		Tags []string `yaml:"tags"`
	} `yaml:"with"`
}

func NewTasks() *Tasks {
	return &Tasks{}
}
