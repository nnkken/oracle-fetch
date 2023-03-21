package main

import (
	"github.com/nnkken/oracle-fetch/cmd"
	"go.uber.org/zap"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		zap.S().Panicw("command failed with error", "error", err)
	}
}
