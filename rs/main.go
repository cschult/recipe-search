package rs

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func Search() {
	cmdName := "recoll"
	cmdArgs := []string{"-c", "/home/schulle/.config/recoll", "-t", "-b", "dir:/home/schulle/ownCloud/rezepte"}

	fmt.Print("Enter search: ")
	searchTerm := bufio.NewScanner(os.Stdin)
	searchTerm.Scan()
	search := searchTerm.Text()
	// fmt.Printf("%s\n", search)
	cmdArgs = append(cmdArgs, search)

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
	for i, v := range result {
		fmt.Println(i, v)
	}
}