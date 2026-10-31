package main

import "fmt"

func main() {
	keys := []string{"abc"}
	runKeys := keys

	runKeys = append(runKeys, "bcd")

	fmt.Println("keys:", keys)
	fmt.Println("runKeys:", runKeys)
}
