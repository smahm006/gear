package tasks

type Tasks struct {
	Tag  string
	Name string        `yaml:"name"`
	With PreTaskLogic  `yaml:"with"`
	And  PostTaskLogic `yaml:"and"`
}

func NewTasks() *Tasks {
	return &Tasks{}
}
