package main

import (
	"github.com/seaung/nox/pkg/cmd"
	"github.com/seaung/nox/pkg/utils"
)

func main() {
	utils.InitConsole()
	cmd.Execute()
}
