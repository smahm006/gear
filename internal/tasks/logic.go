package tasks

type PreTaskLogic struct {
	Conditions string   `yaml:"conditions"`
	Sudo       bool     `yaml:"sudo"`
	Items      []string `yaml:"items"`
	Tags       []string `yaml:"tags"`
	DelegateTo string   `yaml:"delegate_to"`
}

type PostTaskLogic struct {
	IgnoreErrors string `yaml:"ignore_erros"`
	Retry        int    `yaml:"retry"`
	Delay        int    `yaml:"delay"`
	Notify       string `yaml:"notipy"`
}
