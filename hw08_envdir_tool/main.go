package main

import (
	"os"
)

func main() {
	directoryName := os.Args[1]
	env, err := ReadDir(directoryName)
	if err != nil {
		return
	}
	RunCmd(os.Args[2:], env)
}
