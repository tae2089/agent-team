package main

import (
	"os"

	"github.com/tae2089/agent-team/internal/cli"
	"github.com/tae2089/agent-team/internal/output"
)

func main() {
	if err := cli.NewRoot().Execute(); err != nil {
		output.WriteError(os.Stdout, 0, err)
		os.Exit(1)
	}
}
