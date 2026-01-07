package main

import (
	"os"

	"github.com/JamesIOmete/platformctl/internal/cli"
)

func main() {
	code := cli.Run(os.Args[1:])
	if code != 0 {
		os.Exit(code)
	}
}
