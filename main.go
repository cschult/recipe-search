package main

import (
	"devmem.de/srv/git/recipe-search/internal/h"
	"devmem.de/srv/git/recipe-search/internal/rs"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strconv"
)

// todo: message when search has no result: sorry and show new search dialog prompt, do not print helpline
// todo: ensure that only txt files are concatenated
// todo: was passiert, wenn Corinna eine Datei editiert?
// todo: remember state l or s


// main is the text user interface of this program.
// It starts a loop for interacting with the user.
func main() {

	resultPathFile, resultFile := rs.Search()
	rs.ViewResult(resultFile)
	helpLine := color.CyanString("h help | q quit | n new search | l long |" +
		" s short | p print | e edit | 1 = show file #1 | 2 = ...")

	for true {
		fmt.Println(helpLine)
		fmt.Printf("Enter a key: ")
		key := rs.Input()

		switch key {
		case "h":
			h.Help()
		case "q":	// quit
			os.Exit(0)
		case "n":	// new search
			// clear the list of command line args so that rs.args()
			// asks user for new search term
			myName := os.Args[0]
			os.Args = []string{myName}
			resultPathFile, resultFile = rs.Search()
			rs.ViewResult(resultFile)
		case "l":	// path and filename
			rs.ViewResult(resultPathFile)
		case "s":	// only filename
			rs.ViewResult(resultFile)
		case "e":
			rs.EditFile(resultFile, resultPathFile)
		case "p":	// send file to printer
			rs.Print(resultFile, resultPathFile)
		case "":	// ENTER, print help line
			rs.ViewResult(resultFile)
			color.Yellow("enter a valid key or a file number\n\n")
		default:
			i, err := strconv.Atoi(key)
			if err == nil {
				if i >= 1 && i <= len(resultPathFile) {
					rs.ConcatFile(resultFile, resultPathFile, i-1)
				} else {
					rs.ViewResult(resultFile)
					color.Yellow("not a valid file number\n\n")
				}
			} else {
				rs.ViewResult(resultFile)
				color.Yellow("invalid key\n\n")
			}
		}
	}
}
