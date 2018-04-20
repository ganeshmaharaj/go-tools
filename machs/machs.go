package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"gopkg.in/yaml.v2"
)

type mach_info struct {
	Ip   string `yaml:"ip"`
	User string `yaml:"user,omitempty"`
}

func main() {
	root_dir := filepath.Join(os.Getenv("HOME"), "bin/servers")
	// Map file with machine name, IP and User
	mymachs := make(map[string]mach_info)
	// Created a keyarray as i am unable to index a map and use the int value
	keyarr := []string{}
	// Get User input
	buffin := bufio.NewScanner(os.Stdin)

	// Read all yaml/yml files and construct the map with all the details
	files, _ := filepath.Glob(root_dir + "/*.y*ml")
	for _, file := range files {
		yamlfile, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("Unable to read file")
			return
		}

		err = yaml.Unmarshal([]byte(yamlfile), &mymachs)
		if err != nil {
			fmt.Printf("Unable to unmarshal file, %s", file)
			return
		}
	}

	// Create the array with map key values and show that to the user. Use the
	// return index value to get the key which is then passed to map to get Ip
	// and User
	for key, _ := range mymachs {
		keyarr = append(keyarr, key)
	}
	fmt.Println("Index	Machine")
	fmt.Println("===============")
	for idx, val := range keyarr {
		fmt.Printf("%d	%s\n", idx, val)
	}
	fmt.Print("Machine >> ")
	buffin.Scan()
	userval, err := strconv.Atoi(buffin.Text())
	if err != nil {
		fmt.Println("Unable to get convert user input")
		return
	}

	// Construct args of ssh command.
	// Use syscall.Exec to replace the running process itself
	mycmd := []string{"-X"}
	if mymachs[keyarr[userval]].User != "" {
		mycmd = append(mycmd, mymachs[keyarr[userval]].User+"@"+mymachs[keyarr[userval]].Ip)
	} else {
		mycmd = append(mycmd, mymachs[keyarr[userval]].Ip)
	}
	sshbin, err := exec.LookPath("ssh")
	env := os.Environ()

	err = syscall.Exec(sshbin, mycmd, env)
	if err != nil {
		panic(err)
	}
}
