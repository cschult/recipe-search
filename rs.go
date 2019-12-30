package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// docker build current directory
	cmdName := "recoll"
	cmdArgs := []string{"-c", "/home/schulle/.config/recoll", "-t", "-b", "dir:/home/schulle/ownCloud/rezepte"}
	// cmdArgs := []string{"build", "."}

	fmt.Print("Enter search: ")
	searchTerm := bufio.NewScanner(os.Stdin)
	searchTerm.Scan()
	search := searchTerm.Text()
	cmdArgs = append(cmdArgs, search)

	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
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
}
