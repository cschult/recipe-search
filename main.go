package main

// todo: show pdf files ('pdftotext -' or external viewer)
// todo: add printing

import (
	"bufio"
	"devmem.de/srv/git/recipe-search/internal"
	"fmt"
	"os"
	"strconv"
)


func main() {
	resultPathFile, resultFile := rs.Search()
	rs.ViewResult(resultFile)

	for true {
		fmt.Println("q: quit; n: search; l: long; s: short; 1 = show file #1; 2 = ...")
		fmt.Printf("Enter a key: ")
		// input := bufio.NewScanner(os.Stdin)
		// input.Scan()
		// key := input.Text()
		key := rs.Input()

		switch key {
		case "q":	// quit
			os.Exit(0)
		case "n":	// new search
			// empty list of command line args so that rs.args()
			// asks user for new search term
			myName := os.Args[0]
			os.Args = []string{myName}
			resultPathFile, resultFile = rs.Search()
			rs.ViewResult(resultFile)
		case "l":
			rs.ViewResult(resultPathFile)
		case "s":
			rs.ViewResult(resultFile)
		// case "p":	// print
		// 	rs.Print(resultPathFile, 1)
		case "":	// ENTER
			fmt.Println("q = quit, 1 = show file #1, ...")
		default:
			i, err := strconv.Atoi(key)
			if err == nil {
				if i >= 1 && i <= len(resultPathFile) {
					rs.FileConcat(resultPathFile, i-1)
					fmt.Println("press ENTER to continue")
					cont := bufio.NewScanner(os.Stdin)
					cont.Scan()
					rs.ViewResult(resultFile)

				} else {
					fmt.Println("not a valid file number")
				}
			}
		}
	}
}
