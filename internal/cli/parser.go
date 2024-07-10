package cli

import (
	"flag"
	"os"
	"strings"
	"text/template"

	"github.com/smahm006/gear/internal/utils"
)

type CliParser struct {
	Help          bool
	Version       bool
	PlaybookPath  string
	InventoryPath string
	RolePaths     []string
	Verbosity     int
	ExtraVars     []string
	Tags          string
	Limit         string
}

func NewCliParser() *CliParser {
	return &CliParser{
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

func parseVerbosity(args []string) (int, []string) {
	verbosity_count := 0
	remaining_args := []string{}
	for _, arg := range args {
		if strings.HasPrefix(arg, "-v") {
			verbosity_count += len(arg) - 1
		} else {
			remaining_args = append(remaining_args, arg)
		}
	}
	return verbosity_count, remaining_args
}

func (p *CliParser) Parse() error {
	flag.Usage = func() {
		utils.CheckErr(ShowUsage())
	}
	// Need to parse verbosity first
	verbosity, remaining_args := parseVerbosity(os.Args[1:])
	os.Args = append([]string{os.Args[0]}, remaining_args...)
	p.Verbosity = verbosity
	// Parse the remaining flags
	flag.BoolVar(&p.Help, "h", false, "show help")
	flag.BoolVar(&p.Help, "help", false, "show help")
	flag.BoolVar(&p.Version, "version", false, "show version")
	flag.StringVar(&p.InventoryPath, "i", "", "path to inventory")
	flag.StringVar(&p.InventoryPath, "inventory", "", "path to inventory")
	flag.StringVar(&p.PlaybookPath, "p", "", "paths to playbooks")
	flag.StringVar(&p.PlaybookPath, "playbook", "", "paths to playbooks")
	// flag.StringVar(&p., "u", "", "username for ssh connections")
	// flag.StringVar(&p., "user", "", "username for ssh connection")
	flag.StringVar(&p.Limit, "limit", "", "group/hosts limit")
	flag.StringVar(&p.Limit, "l", "", "group/hosts limit")
	flag.StringVar(&p.Tags, "tags", "", "tasks limit")
	flag.StringVar(&p.Tags, "t", "", "tasks limit")
	flag.Parse()
	return nil
}
