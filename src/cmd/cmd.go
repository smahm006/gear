package cmd

func NewGearCommand() (*CliParser, error) {
	var cmd CliParser
	cmd = *NewCliParser()
	cmd.Parse()
	return &cmd, nil
}
