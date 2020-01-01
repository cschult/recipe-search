package rs

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Input() string {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	inputText := input.Text()
	return inputText
}

// add search words to list of arguments
// get from command line or ask at runtime
func args(cmdArgs []string) []string {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		fmt.Print("Enter search: ")
		// searchTerm := bufio.NewScanner(os.Stdin)
		// searchTerm.Scan()
		// search := searchTerm.Text()
		// fmt.Printf("%s\n", search)
		search := Input()
		cmdArgs = append(cmdArgs, search)
		return cmdArgs
	} else {
		cmdArgs = append(cmdArgs, argsWithoutProg...)
		return cmdArgs
	}
}


func Search() ([]string, []string) {
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
	var resultPathFile []string
	var resultFile []string
	go func() {
		for scanner.Scan() {
			resultPathFile = append(resultPathFile, scanner.Text())
			}
			for _, f := range resultPathFile {
				// resultFile = append(resultFile, strings.TrimPrefix(f, "file://"))
				resultFile = append(resultFile, filepath.Base(f))
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
	return resultPathFile, resultFile
}


func ViewResult(result []string)  {
	for i, v := range result {
		fmt.Println(i+1, v)
	}
	fmt.Println()
}


func FileConcat(res []string, i int)  {

	f := res[i]
	f = strings.TrimPrefix(f, "file://")
	fmt.Println("file:", f)
	file, err := os.Open(f)

	if err != nil {
		fmt.Fprintln(os.Stderr, "failed opening file: %s\n", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	fmt.Println()
}

func Print(res []string, i int) {
	// 	prints recipe to printer
	// - need file number to select file from slice of files
	// check $PRINTER
	// or list all available printers with lpstat -a
	// ask user for printer selection
	// print file
}
