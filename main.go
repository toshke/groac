package main

import (
	"os"

	cmd "github.com/toshke/groac/cmd"
)

var GROAC_VERSION = "0.0.1"

func init() {
	os.Setenv("GROAC_VERSION", GROAC_VERSION)
}

func main() {
	cmd.Execute()
}
