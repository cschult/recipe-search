package rs

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/akutz/sortfold"
	"github.com/fatih/color"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// =================
// F U N C T I O N S
// =================

// Search calls external program 'recoll', a document indexer, with arguments
// and search term and returns two slices with the list of matching files.
// Slice 'resultPathFile' has the file names with path, resultFile only the
// file names.
func Search(searcher string, searcherArgs []string) ([]string, []string) {

	// call args to get the command line args and create the Cmd struct 'cmd'
	// and 'reader' as a handle for stdout of external program.
	searcherArgs = args(searcherArgs)
	cmd := exec.Command(searcher, searcherArgs...)
	reader, err := cmd.StdoutPipe()
	if err != nil {
		PrtErr("Error creating StdoutPipe for Cmd:", err)
		os.Exit(1)
	}

	// read stdout line wise and return two slices with the file names
	// (with and without path)
	scanner := bufio.NewScanner(reader)
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
		fmt.Println(searcher, searcherArgs)
		PrtErr("Error starting Cmd:", err)
		os.Exit(1)
	}

	// wait for external program to finish and error handling
	err = cmd.Wait()
	if err != nil {
		PrtErr("Error waiting for Cmd:", err)
		os.Exit(1)
	}

	mySort(resultPathFile)
	mySort(resultFile)

	return resultPathFile, resultFile
}

// ConcatFile reads given file and prints it on screen line by line.
// Then waiting for ENTER key to show list of files again.
func ConcatFile(resultPathFile []string, resultFile []string, i int) error {

	f := resultPathFile[i]
	f = strings.TrimPrefix(f, "file://")
	fmt.Println("file:", f)
	file, err := os.Open(f)

	if err != nil {
		// PrtErr("failed to open file:", err)
		return err
	}

	defer FileClose(file)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	fmt.Printf("\npress ENTER to continue")
	cont := bufio.NewScanner(os.Stdin)
	cont.Scan()
	return err
}

// EditFile opens external editor to edit a recipe
func EditFile(editor string, resultPathFile []string) error {

	var file string
	fmt.Printf("enter number of file to edit: ")
	k := Input()

	if k == "" {
		color.Yellow("pressed ENTER, returning\n\n")
		return nil
	}

	i, err := strconv.Atoi(k)
	if err == nil {
		if i >= 1 && i <= len(resultPathFile) {
			file = resultPathFile[i-1 ]
			file = strings.TrimPrefix(file, "file://")
		} else {
			color.Yellow("no file with that number, returning\n\n")
			return nil
		}
	} else {
		color.Yellow("not a number, returning\n\n")
		return nil
	}

	editorArgs := []string{file}

	cmd := exec.Command(editor, editorArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		// PrtErr("Error starting editor:", err)
		return errors.New("error starting editor")
	}

	err = cmd.Wait()
	if err != nil {
		// PrtErr("Error waiting for editor:", err)
		return errors.New("error waiting for editor")
	}
	return nil
}

// 	Print sends recipe to printer
func Print(cfg map[string]string, resultPathFile []string) error {

	prntcmd := cfg["prntcmd"]
	converter := cfg["converter"]
	prntcmdArgs := cfg["prntcmdArgs"]
	converterArgs := cfg["converterArgs"]
	printer := cfg["printer"]
	prntduplex := cfg["prntduplex"]
	prntcolor := cfg["prntcolor"]

	// fixme: same code here as in edit function
	var file string
	fmt.Printf("enter number of file to print: ")
	k := Input()

	if k == "" {
		color.Yellow("pressed ENTER, returning\n\n")
		return nil
	}

	i, err := strconv.Atoi(k)
	if err == nil {
		if i >= 1 && i <= len(resultPathFile) {
			file = resultPathFile[i-1 ]
			file = strings.TrimPrefix(file, "file://")
		} else {
			// color.Yellow("no file with that number, returning\n\n")
			return errors.New("no file with that number, returning")
		}
	} else {
		color.Yellow("not a number, returning\n\n")
		return nil
	}

	// test if file exists and is readable
	var f *os.File
	f, err = os.Open(file)
	if err != nil {
		// PrtErr("failed to open file:", err)
		return err
	}
	FileClose(f)

	// firstArgs := []string{txtConverterArgs, file}
	firstArgs := []string{converterArgs, file}
	first := exec.Command(converter, firstArgs...)
	secondArgs := []string{prntcmdArgs, printer, prntduplex, prntcolor}
	second := exec.Command(prntcmd, secondArgs...)

	reader, writer := io.Pipe()
	first.Stdout = writer
	second.Stdin = reader

	err = first.Start()
	if err != nil {
		// PrtErr("Error starting paps:", err)
		return errors.New("error starting paps")
	}

	err = second.Start()
	if err != nil {
		// PrtErr( "Error starting lpr:", err)
		return errors.New("error starting lpr")
	}

	err = first.Wait()
	if err != nil {
		// PrtErr("Error waiting for paps:", err)
		return errors.New("error waiting for paps")
	}

	writer.Close()
	err = second.Wait()
	if err != nil {
		// PrtErr("Error waiting for lpr:", err)
		return errors.New("error waiting for lpr")
	}

	fmt.Printf("printed file %s\n", file)
	return nil
}

// ===============================
// H E L P E R   F U N C T I O N S
// ===============================

// ViewResult print the matching files on screen
func ViewResult(resultpath []string, result []string, l bool) {
	fmt.Println()
	if l == true {
		for i, v := range resultpath {
			fmt.Println(i+1, v)
		}
	} else {
		for i, v := range result {
			fmt.Println(i+1, v)
		}
	}
	fmt.Println()
}

// Input collects user input from keyboard and returns it as string.
func Input() string {
	in := bufio.NewScanner(os.Stdin)
	in.Scan()
	s := in.Text()
	return s
}

// mySort sorts slice of strings case insensitive
func mySort(a []string) {
	sort.SliceStable(a, func(i, j int) bool {
		return sortfold.CompareFold(a[i], a[j]) < 0
	})
}

// FileClose closes a file and exits if error occurs
func FileClose(f *os.File) {
	err := f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// PrtErr print to stderr with color red,
// surrounds string with newline
func PrtErr(s string, err error) {
	red := color.New(color.FgHiRed).FprintfFunc()
	red(os.Stderr, "%s %s\n\n", s, err)
}

// args adds search words to list of arguments.
// It get's them from command line or asks the user,
// if no command line args where given.
// fixme: split in two funcs, so main.go/main must not clear the list of cmd line args when a new search is started with "n".
func args(searcherArgs []string) []string {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Print("\nEnter search: ")
		search := Input()
		searcherArgs = append(searcherArgs, search)
		return searcherArgs
	} else {
		searcherArgs = append(searcherArgs, args...)
		return searcherArgs
	}
}
