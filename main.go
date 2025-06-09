package main

import (
	"github.com/Smartling/smartling-cli/cmd"
)

func main() {
	rootCmd, err := cmd.NewRootCmd()
	if err != nil {
		panic(err)
	}
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
