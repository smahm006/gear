package cmd

import (
	"flag"
	"os"
	"text/template"

	"github.com/smahm006/gear/lib/logger"
)

type CLIParser struct {
	Help          bool
	Version       bool
	PlaybookPaths []string
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

func ShowHelp() error {
	t := template.New("help")
	template.Must(t.Parse("usage:\n"))
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
		logger.CheckErr(ShowHelp())
	}
	flag.BoolVar(&p.Help, "help", false, "show help")
	flag.BoolVar(&p.Version, "version", false, "show version")
	flag.StringVar(&p.InventoryPath, "i", "", "show version")
	flag.Parse()
	return nil
}
