package rs

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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

func Search() []string {
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
	return result
}

func ViewResult(result []string)  {
	for i, v := range result {
		fmt.Println(i+1, v)
	}
}

/*
enter file number to view file or n for new search or q to quit
 */

/* open file + print it out */

func FileConcat(res []string, i int)  {

	f := res[i]
	f = strings.TrimPrefix(f, "file://")
	fmt.Println("opening:", f)
	file, err := os.Open(f)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}

}
