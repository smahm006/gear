package cmd

import (
	"flag"
	"os"
	"text/template"

	"github.com/smahm006/gear/src/utils"
)

type CLIParser struct {
	Help          bool
	Version       bool
	PlaybookPath  string
	InventoryPath string
	RolePaths     []string
	Verbosity     uint32
	Tags          []string
	ExtraVars     []string
}

func NewCliParser() *CLIParser {
	return &CLIParser{
		Help:      false,
		Version:   false,
		Verbosity: 0,
	}
}

func ShowUsage() error {
	const usage = `Usage:
    gear [--inventory] [--playbook] [--user]
Options:
    -i, --inventory             Path to the inventory file.
    -p, --playbook              Path to the playbook file.
    -u, --user USER             Username used for ssh connection.
`
	t := template.New("usage")
	template.Must(t.Parse(usage))
	if err := t.Execute(os.Stdout, nil); err != nil {
		return err
	}
	return nil
}

func ShowVersion() error {
	v := struct {
		Version string
		Build   string
	}{
		Version: VERSION,
		Build:   BUILD,
	}
	version_template := "Gear version {{.Version}}, build {{.Build}}\n"
	t := template.New("version")
	template.Must(t.Parse(version_template))
	if err := t.Execute(os.Stdout, v); err != nil {
		return err
	}
	return nil
}

func (p *CLIParser) Parse() error {
	flag.Usage = func() {
		utils.CheckErr(ShowUsage())
	}
	flag.BoolVar(&p.Help, "h", false, "show help")
	flag.BoolVar(&p.Help, "help", false, "show help")
	flag.BoolVar(&p.Version, "v", false, "show version")
	flag.BoolVar(&p.Version, "version", false, "show version")
	flag.StringVar(&p.InventoryPath, "i", "", "path to inventory")
	flag.StringVar(&p.InventoryPath, "inventory", "", "path to inventory")
	flag.StringVar(&p.PlaybookPath, "p", "", "paths to playbooks")
	flag.StringVar(&p.PlaybookPath, "playbook", "", "paths to playbooks")
	flag.StringVar(&p.PlaybookPath, "u", "", "username for ssh connections")
	flag.StringVar(&p.PlaybookPath, "user", "", "username for ssh connection")
	flag.Parse()
	return nil
}
