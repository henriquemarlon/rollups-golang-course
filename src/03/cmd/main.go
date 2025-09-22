package main

import (
	"os"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/cmd/root"
)

func main() {
	err := root.Cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
