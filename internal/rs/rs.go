package rs

import (
	"bufio"
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


// global config vars
var (
	searcher		 = "recoll"
	editor           = "nvim"
	txtConverter     = "paps"
	lpr              = "lpr"
	lprArgs          = "-P"
	txtConverterArgs = "--font=Monospace 10"
	printer          = "GraustufenNormalDuplex"
	printDuplex      = "-o Duplex=DuplexNoTumble"
	printcolor       = "-o BRMonoColor=Mono"
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

	// cmdName := "recoll"
	// arguments to call recoll for my recipe collection
	searcherArgs := []string{"-t", "-c", "/home/schulle/.config/recoll", "-t", "-b", "dir:/home/schulle/ownCloud/rezepte"}

	// call args to get the command line args and create the Cmd struct 'cmd'
	// and 'reader' as a handle for stdout of external program.
	searcherArgs = args(searcherArgs)
	cmd := exec.Command(searcher, searcherArgs...)
	reader, err := cmd.StdoutPipe()
	if err != nil {
		prtErr("Error creating StdoutPipe for Cmd:", err)
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
		prtErr("Error starting Cmd:", err)
		os.Exit(1)
	}

	// wait for external program to finish and error handling
	err = cmd.Wait()
	if err != nil {
		prtErr("Error waiting for Cmd:", err)
		os.Exit(1)
	}

	mySort(resultPathFile)
	mySort(resultFile)

	return resultPathFile, resultFile
}


// mySort sorts slice of strings case insensitive
func mySort (a []string) {
	sort.SliceStable(a, func(i, j int) bool {
		return sortfold.CompareFold(a[i], a[j]) < 0
	})
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


// fileClose closes a file and exits if error occurs
func fileClose(f *os.File)  {
	err := f.Close()
	if err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
}


// ConcatFile reads given file and prints it on screen line by line.
// Then waiting for ENTER key to show list of files again.
func ConcatFile(resultFile []string, resultPathFile []string, i int)  {

	f := resultPathFile[i]
	f = strings.TrimPrefix(f, "file://")
	fmt.Println("file:", f)
	file, err := os.Open(f)

	if err != nil {
		prtErr("failed to open file:", err)
	}

	defer fileClose(file)

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


// EditFile opens external editor to edit a recipe
func EditFile(resultFile []string, resultPathFile []string)  {

	var file string
	fmt.Printf("enter number of file to edit: ")
	k := Input()

	if k == "" {
		ViewResult(resultFile)
		color.Yellow("pressed ENTER, returning\n\n")
		return
	}

	i, err:= strconv.Atoi(k)
	if err == nil {
		if i >= 1 && i <= len(resultPathFile) {
			file = resultPathFile[i-1 ]
			file = strings.TrimPrefix(file, "file://")
		} else {
			ViewResult(resultFile)
			color.Yellow("no file with that number, returning\n\n")
			return
		}
	} else {
		ViewResult(resultFile)
		color.Yellow("not a number, returning\n\n")
		return
	}

	editorArgs := []string{file}

	cmd := exec.Command(editor, editorArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		ViewResult(resultFile)
		prtErr("Error starting editor:", err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		ViewResult(resultFile)
		prtErr("Error waiting for editor:", err)
		return
	}
	ViewResult(resultFile)
}


// prtErr print to stderr with color red,
// surrounds string with newline
func prtErr(s string, err error) {
	red := color.New(color.FgHiRed).FprintfFunc()
	red(os.Stderr, "%s %s\n\n", s, err)
}


// 	Print prints recipe to printer
func Print(resultFile []string, resultPathFile []string) {

	// fixme: same code here as in edit function
	var file string
	fmt.Printf("enter number of file to print: ")
	k := Input()

	if k == "" {
		ViewResult(resultFile)
		color.Yellow("pressed ENTER, returning\n\n")
		return
	}

	i, err:= strconv.Atoi(k)
	if err == nil {
		if i >= 1 && i <= len(resultPathFile) {
			file = resultPathFile[i-1 ]
			file = strings.TrimPrefix(file, "file://")
		} else {
			ViewResult(resultFile)
			color.Yellow("no file with that number, returning\n\n")
			return
		}
	} else {
		ViewResult(resultFile)
		color.Yellow("not a number, returning\n\n")
		return
	}

	firstArgs := []string{txtConverterArgs, file}
	first := exec.Command(txtConverter, firstArgs...)
	secondArgs := []string{lprArgs, printer, printDuplex, printcolor}
	second := exec.Command(lpr, secondArgs...)

	reader, writer := io.Pipe()
	first.Stdout = writer
	second.Stdin = reader

    err = first.Start()
	if err != nil {
		ViewResult(resultFile)
		prtErr("Error starting paps:", err)
		return
	}

    err = second.Start()
	if err != nil {
		ViewResult(resultFile)
		prtErr("Error starting lpr:", err)
		return
	}

	err = first.Wait()
	if err != nil {
		ViewResult(resultFile)
		prtErr("Error waiting for paps:", err)
		return
	}

	writer.Close()
    err = second.Wait()
	if err != nil {
		ViewResult(resultFile)
		prtErr("Error waiting for lpr:", err)
		return
	}

	fmt.Printf("printed file %s\n", file)
	ViewResult(resultFile)
}
