package state

import (
	"github.com/smahm006/gear/internal/cli"
	"github.com/smahm006/gear/internal/inventory"
)

// State - The particular condition that someone or something is in at a specific time. ‘the state of the company's finances’
type RunState struct {
	ParsedFlags *cli.CliParser
	Status      *RunStatus
	Inventory   *inventory.Inventory
}

func NewRunState(cli *cli.CliParser, i *inventory.Inventory) *RunState {
	return &RunState{
		ParsedFlags: cli,
		Inventory:   i,
	}
}
