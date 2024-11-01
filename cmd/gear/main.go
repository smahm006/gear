package main

import (
	"fmt"
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
	cmd, err := cli.NewGearCommand()
	if err != nil {
		return err
	}
	if len(os.Args) == 1 || os.Args[1] == "help" || cmd.Help {
		cli.ShowUsage()
		return nil
	}
	if os.Args[1] == "version" || cmd.Version {
		cli.ShowVersion()
		return nil
	}
	if len(cmd.InventoryPath) == 0 {
		return fmt.Errorf("no inventory path provided")
	}
	if len(cmd.PlaybookPath) == 0 {
		return fmt.Errorf("no playbook path provided")
	}
	i := inventory.NewInventory()
	if err = i.LoadInventory(cmd.InventoryPath); err != nil {
		return err
	}
	p := playbook.NewPlaybook()
	if err = p.LoadPlaybook(cmd, i); err != nil {
		return err
	}
	if err = p.RunPlaybook(cmd, i); err != nil {
		return err
	}
	return nil
}
