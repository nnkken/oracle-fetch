package main

import (
	"github.com/nnkken/oracle-fetch/cmd"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		// TODO: proper logging
		panic(err)
	}
}
