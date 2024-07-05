package utils

import (
	"encoding/json"
	"fmt"
)

func PrintMap(i interface{}) {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print(string(b))
}
