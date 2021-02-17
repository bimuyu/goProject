package Demo

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func CmdDemo(cmd string) {
	command := exec.Command(cmd)
	pipe, err := command.StdoutPipe()
	for err != nil {
		fmt.Println(err, 1)
		return
	}
	if err := command.Start(); err != nil {
		fmt.Println(err, 2)
		return
	}
	all, err := ioutil.ReadAll(pipe)
	pipe.Close()
	if err := command.Wait(); err != nil {
		fmt.Println(err, 3)
		return
	}
	fmt.Println(string(all))

}
