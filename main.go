package main

import (
	"os"

	"github.com/smahm006/gear/lib/logger"
	"github.com/smahm006/gear/src/cmd"
	"github.com/smahm006/gear/src/inventory"
	"github.com/smahm006/gear/src/playbooks"
)

func main() {
	logger.CheckErr(entrypoint())
}

func entrypoint() error {
	cli, err := cmd.NewGearCommand()
	if err != nil {
		return err
	}
	if len(os.Args) == 1 || cli.Help {
		cmd.ShowHelp()
		return nil
	}
	if os.Args[1] == "version" || cli.Version {
		cmd.ShowVersion()
		return nil
	}
	i := inventory.NewInventory()
	if err = i.LoadInventory(cli.InventoryPath); err != nil {
		return err
	}
	p := playbooks.NewPlaybook()
	if err = p.LoadPlaybook(cli.PlaybookPath); err != nil {
		return err
	}
	return nil
}
