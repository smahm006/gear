package main

import (
	"os"

	"github.com/smahm006/gear/lib/logger"
	"github.com/smahm006/gear/src/cmd"
	"github.com/smahm006/gear/src/inventory"
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
	inventory := inventory.NewInventory()
	return nil
}
