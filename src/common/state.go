package common

import (
	"github.com/smahm006/gear/src/cmd"
	"github.com/smahm006/gear/src/inventory"
)

// State - The particular condition that someone or something is in at a specific time. ‘the state of the company's finances’
type RunState struct {
	ParsedFlags  *cmd.CliParser
	Inventory    *inventory.Inventory
	PlaybookPath string
	Variables    map[string]interface{}
	Tags         []string
}

func NewRunState(cli *cmd.CliParser, i *inventory.Inventory) *RunState {
	return &RunState{
		ParsedFlags: cli,
		Inventory:   i,
	}
}
