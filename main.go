package main

import (
	"github.com/vincentdaniel/subify/cmd"
	"github.com/vincentdaniel/subify/common/utils"
)

func main() {
	utils.InitLoggingConf()
	cmd.Execute()
}
