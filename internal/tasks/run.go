package tasks

import (
	"fmt"
	"sync"

	"github.com/smahm006/gear/internal/inventory"
	"github.com/smahm006/gear/internal/playbook/state"
	"github.com/smahm006/gear/internal/tasks/exchange"
)

func (t *Task) RunTask(state *state.RunState) {
	resp_chan := make(chan *exchange.TaskResponse)
	var wg_command sync.WaitGroup
	var wg_processing sync.WaitGroup
	for _, host := range state.Status.Hosts {
		fmt.Printf("running task %s\n", t.Name)
		wg_command.Add(1)
		go func(host *inventory.Host, resp_chan chan *exchange.TaskResponse) {
			defer wg_command.Done()
			connection, err := state.Status.ConnectionCache.GetConnection(host)
			if err != nil {
				fmt.Println(err)
			}
			connection.Connect()
			request := t.Module.Query()
			response := t.Module.Run(connection, request, t.With, t.And)
			if err != nil {
				fmt.Println(err)
			}
			resp_chan <- response
		}(host, resp_chan)
	}
	wg_processing.Add(1)
	go func() {
		defer wg_processing.Done()
		for resp := range resp_chan {
			fmt.Println(resp.CommandResult.Out)
		}
	}()
	wg_command.Wait()
	close(resp_chan)
	wg_processing.Wait()
}
