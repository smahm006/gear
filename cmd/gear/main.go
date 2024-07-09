package main

import (
	"os"

	"github.com/smahm006/gear/internal/cli"
	"github.com/smahm006/gear/internal/inventory"
	"github.com/smahm006/gear/internal/playbook"
	"github.com/smahm006/gear/internal/utils"
)

func main() {
	utils.CheckErr(entrypoint())
}

func entrypoint() error {
	cli, err := cli.NewGearCommand()
	if err != nil {
		return err
	}
	if len(os.Args) == 1 || os.Args[1] == "help" || cli.Help {
		cli.ShowUsage()
		return nil
	}
	if os.Args[1] == "version" || cli.Version {
		cli.ShowVersion()
		return nil
	}
	i := inventory.NewInventory()
	if err = i.LoadInventory(cli.InventoryPath); err != nil {
		return err
	}
	p := playbook.NewPlaybook()
	if err = p.LoadPlaybook(cli, i); err != nil {
		return err
	}
	if err = p.RunPlaybook(cli, i); err != nil {
		return err
	}
	return nil
}
