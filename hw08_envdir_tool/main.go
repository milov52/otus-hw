package main

import (
	"fmt"
	"os"
)

func main() {
	directoryName := os.Args[1]
	env, err := ReadDir(directoryName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
		os.Exit(1)
	}
	returnCode := RunCmd(os.Args[2:], env)
	if returnCode != 0 {
		fmt.Fprintf(os.Stderr, "Command exited with code: %d\n", returnCode)
		os.Exit(returnCode)
	}
}
