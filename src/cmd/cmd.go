package cmd

func NewGearCommand() (CLIParser, error) {
	var cmd CLIParser
	cmd = *NewCliParser()
	cmd.Parse()
	return cmd, nil
}
