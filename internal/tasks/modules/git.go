package modules

type GitModule struct {
	Repository  string `yaml:"repo"`
	Branch      string `yaml:"branch"`
	FailedWhen  string `yaml:"failed_when"`
	ChangedWhen string `yaml:"changed_when"`
}
