package playbook

import "github.com/smahm006/gear/internal/playbook/state"

func CleanUpPlay(run_status *state.RunStatus) {
	for _, connection := range run_status.ConnectionCache.Connections {
		connection.Close()
	}
}
