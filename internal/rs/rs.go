package rs

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)


// Input collects user input from keyboard and returns it as string.
func Input() string {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	inputText := input.Text()
	return inputText
}


// args adds search words to list of arguments.
// It get's them from command line or asks the user,
// if no command line args where given.
// fixme: split in two funcs, so main.go/main must not clear the list of cmd line args when a new search is started with "n".
func args(cmdArgs []string) []string {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		fmt.Print("\nEnter search: ")
		search := Input()
		cmdArgs = append(cmdArgs, search)
		return cmdArgs
	} else {
		cmdArgs = append(cmdArgs, argsWithoutProg...)
		return cmdArgs
	}
}


// Search calls external program 'recoll', a document indexer, with arguments
// and search term and returns two slices with the list of matching files.
// Slice 'resultPathFile' has the file names with path, resultFile only the
// file names.
func Search() ([]string, []string) {
	cmdName := "recoll"
	// arguments to call recoll for my recipe collection
	cmdArgs := []string{"-c", "/home/schulle/.config/recoll", "-t", "-b", "dir:/home/schulle/ownCloud/rezepte"}

	// call args to get the command line args and create the Cmd struct 'cmd'
	// and 'cmdReader' as a handle for stdout of external program.
	cmdArgs = args(cmdArgs)
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	// read stdout line wise and return two slices with the file names
	// (with and without path)
	scanner := bufio.NewScanner(cmdReader)
	var resultPathFile []string
	var resultFile []string
	go func() {
		for scanner.Scan() {
			resultPathFile = append(resultPathFile, scanner.Text())
			resultFile = append(resultFile, filepath.Base(scanner.Text()))
		}
	}()

	// start external program and error handling
	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}

	// wait for external program to finish and error handling
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}

	return resultPathFile, resultFile
}


// ViewResult print the matching files on screen, numbered
// beginning with 1
func ViewResult(result []string)  {
	fmt.Println()
	for i, v := range result {
		fmt.Println(i+1, v)
	}
	fmt.Println()
}


// FileConcat reads given file and prints it on screen line by line.
// Then waiting for ENTER key to show list of files again.
func FileConcat(resultFile []string, resultPathFile []string, i int)  {

	f := resultPathFile[i]
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
	fmt.Printf("\npress ENTER to continue")
	cont := bufio.NewScanner(os.Stdin)
	cont.Scan()
	ViewResult(resultFile)
}


func Print(res []string, i int) {
	// 	prints recipe to printer
	// - need file number to select file from slice of files
	// check $PRINTER
	// or list all available printers with lpstat -a
	// ask user for printer selection
	// print file
}
