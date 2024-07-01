package main

import (
	"encoding/json"
	"fmt"
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
	i := inventory.NewInventory()
	if err = i.LoadInventory(cli.InventoryPath); err != nil {
		return err
	}
	b, err := json.MarshalIndent(i, "", "  ")
	fmt.Println(err)
	fmt.Print(string(b))
	return nil
}
