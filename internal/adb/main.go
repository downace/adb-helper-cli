package adb

import (
	"fmt"
	"github.com/ttacon/chalk"
	"os/exec"
)

var Binary = "adb"

func ExecAdb(args ...string) (output string, err error) {
	fmt.Println("Executing", chalk.Blue.Color(Binary), "with args", chalk.Blue.Color(fmt.Sprint(args)))
	cmd := exec.Command(Binary, args...)
	outputBytes, e := cmd.CombinedOutput()

	if e != nil {
		fmt.Println(chalk.Red.Color(e.Error()))
	}
	return string(outputBytes), e
}
