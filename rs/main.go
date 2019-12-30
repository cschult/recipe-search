package rs

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

// add search words to list of arguments
// get from command line or ask at runtime
func args (cmdargs []string) []string {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		fmt.Print("Enter search: ")
		searchTerm := bufio.NewScanner(os.Stdin)
		searchTerm.Scan()
		search := searchTerm.Text()
		// fmt.Printf("%s\n", search)
		cmdargs = append(cmdargs, search)
		return cmdargs
	} else {
		cmdargs = append(cmdargs, argsWithoutProg...)
		return cmdargs
	}
}

func Search() {
	cmdName := "recoll"
	// needed arguments to call recoll for my recipe collection
	cmdArgs := []string{"-c", "/home/schulle/.config/recoll", "-t", "-b", "dir:/home/schulle/ownCloud/rezepte"}

	cmdArgs = args(cmdArgs)
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	var result []string
	go func() {
		for scanner.Scan() {
			result = append(result, scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}
	// todo: make results viewable
	for i, v := range result {
		fmt.Println(i, v)
	}
}
