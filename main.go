package main

import (
	"os"

	"github.com/marintailor/rcstate/cmd"
)

func main() {
	os.Exit(cmd.Run(os.Args[1:]))
}
