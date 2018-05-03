package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	env := os.Environ()
	var cmdout []byte
	var sess []string
	tmuxlsarg := []string{"ls", "-F#S"}
	tmuxconnarg := []string{"-2", "attach", "-t"}
	buffin := bufio.NewScanner(os.Stdin)
	var userval int
	tmuxbin, err := exec.LookPath("tmux")
	if err != nil {
		fmt.Println("Unable to find tmux. Might want to install it")
		return
	}

	cmd := exec.Command(tmuxbin, tmuxlsarg...)
	cmdout, _ = cmd.Output()
	sess = strings.Split(strings.TrimRight(string(cmdout), "\n"), "\n")
	switch len(sess) {
	case 1:
		tmuxconnarg = append(tmuxconnarg, "attach", "-t", sess[0])
	case 0:
		break
	default:
		fmt.Println("Index 	Session")
		fmt.Println("==================")
		for idx, val := range sess {
			fmt.Printf("%d	%s\n", idx, val)
		}
		fmt.Print("Session >> ")
		buffin.Scan()
		userval, err = strconv.Atoi(buffin.Text())
		if err != nil {
			fmt.Println("Unable to parse user input")
			return
		}
		tmuxconnarg = append(tmuxconnarg, "attach", "-t", sess[userval])
	}

	err = syscall.Exec(tmuxbin, tmuxconnarg, env)
	if err != nil {
		panic(err)
	}
}
