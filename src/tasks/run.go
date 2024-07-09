package tasks

import (
	"fmt"

	"github.com/smahm006/gear/src/common"
)

// func (t *Task) RunTask(status *common.RunStatus) {
// 	out_chan := make(chan string)
// 	var wg_command sync.WaitGroup
// 	var wg_processing sync.WaitGroup
// 	for _, host := range status.Hosts {
// 		fmt.Printf("running task on %s\n", host.Name)
// 		wg_command.Add(1)
// 		go func(host *inventory.Host, out_chan chan string) {
// 			defer wg_command.Done()
// 			connection, err := status.ConnectionCache.GetConnection(host)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 			connection.Connect()
// 			out, err := connection.Execute("whoami")
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 			out_chan <- out
// 		}(host, out_chan)
// 	}
// 	wg_processing.Add(1)
// 	go func() {
// 		defer wg_processing.Done()
// 		for o := range out_chan {
// 			fmt.Println(o)
// 		}
// 	}()
// 	wg_command.Wait()
// 	close(out_chan)
// 	wg_processing.Wait()
// }

func (t *Task) RunTask(status *common.RunStatus) {
	for _, host := range status.Hosts {
		connection, err := status.ConnectionCache.GetConnection(host)
		if err != nil {
			fmt.Println(err)
		}
		if err := connection.Connect(); err != nil {
			fmt.Println(err)
			return
		}
		out, err := connection.Execute("echo $CANYOUSEEME")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(out)

	}
}
