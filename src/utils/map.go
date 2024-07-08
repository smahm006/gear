package utils

import (
	"encoding/json"
	"fmt"
)

func MergeMaps[K comparable, V any](dest, src map[K]V) {
	for key, value := range src {
		dest[key] = value
	}
}

func PrintMap(i interface{}) {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print(string(b))
}
