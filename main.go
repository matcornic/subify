package main

import (
	"github.com/matcornic/subify/cmd"
	"github.com/matcornic/subify/common/utils"
)

func main() {
	utils.InitLoggingConf()
	cmd.Execute()
}
