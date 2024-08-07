package tasks

import (
	"fmt"
	"sync"

	"github.com/smahm006/gear/internal/inventory"
	"github.com/smahm006/gear/internal/playbook/state"
	"github.com/smahm006/gear/internal/tasks/exchange"
	"github.com/smahm006/gear/internal/templar"
)

func (t *Task) RunTask(state *state.RunState, item interface{}) error {
	resp_chan := make(chan *exchange.TaskResponse)
	var wg_command sync.WaitGroup
	var wg_processing sync.WaitGroup

	state.Status.Variables["item"] = item
	task_name_parsed, err := templar.GetParsedTemplate(t.Name, state.Status.Variables)
	if err != nil {
		return err
	}
	for _, host := range state.Status.Hosts {
		fmt.Printf("running task %s on host %s with item %s\n", task_name_parsed, host.Name, item)
		wg_command.Add(1)
		go func(host *inventory.Host, resp_chan chan *exchange.TaskResponse) {
			defer wg_command.Done()
			connection, err := state.Status.ConnectionCache.GetConnection(host)
			if err != nil {
				fmt.Println(err)
			}
			err = connection.Connect()
			if err != nil {
				fmt.Println(err)
				return
			}
			request := t.Module.Query()
			response := t.Module.Run(connection, request, state.Status.Variables)
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
	return nil
}
