package common

import (
	"github.com/smahm006/gear/src/cmd"
	"github.com/smahm006/gear/src/inventory"
)

type RunState struct {
	Inevntory *inventory.Inventory
	Status    RunStatus
	Verbosity int
	Variables map[string]interface{}
	Tags      []string
}

func NewRunState(cli *cmd.CliParser, i *inventory.Inventory) *RunState {
	return &RunState{
		Inevntory: i,
		Status:    *NewRunStatus(),
		Verbosity: cli.Verbosity,
	}
}

type RunStatus struct{}

func NewRunStatus() *RunStatus {
	return &RunStatus{}
}
