package main

// todo: show pdf files ('pdftotext -' or external viewer)
// todo: edit recipe in external editor
// todo: add printing

import (
	"devmem.de/srv/git/recipe-search/internal/rs"
	"fmt"
	"os"
	"strconv"
)


// main is the text user interface of this program.
// It starts a loop for interacting with the user.
func main() {
	resultPathFile, resultFile := rs.Search()
	rs.ViewResult(resultFile)
	helpLine := "q: quit; n: search; l: long; s: short; 1 = show file #1; 2 = ..."

	for true {
		fmt.Println(helpLine)
		fmt.Printf("Enter a key: ")
		key := rs.Input()

		switch key {
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
		// case "p":	// send file printer
		// 	rs.Print(resultPathFile, 1)
		case "":	// ENTER, print help line
			fmt.Println(helpLine)
		default:
			i, err := strconv.Atoi(key)
			if err == nil {
				if i >= 1 && i <= len(resultPathFile) {
					rs.FileConcat(resultFile, resultPathFile, i-1)
				} else {
					fmt.Println("not a valid file number")
				}
			}
		}
	}
}
