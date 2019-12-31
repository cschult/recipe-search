package main

// todo: show pdf files ('pdftotext -' or external viewer)
// todo: add printing

import (
	"bufio"
	"devmem.de/srv/git/recipe-search/rs"
	"fmt"
	"os"
	"strconv"
)


func main() {
	result := rs.Search()
	rs.ViewResult(result)

	for true {
		fmt.Println("q = quit; n = new search; 1 = show file #1; 2 = ...")
		fmt.Printf("Enter a key: ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		key := input.Text()

		switch key {
		case "q":
			os.Exit(0)
		case "n":
			// 8
			os.Args = []string{"recipe-search"}
			result = rs.Search()
			rs.ViewResult(result)
		case "":
			fmt.Println("q = quit, 1 = show file #1, ...")
		default:
			i, err := strconv.Atoi(key)
			if err == nil {
				if i >= 1 && i <= len(result) {
					rs.FileConcat(result, i-1)
					fmt.Println("press ENTER to continue")
					cont := bufio.NewScanner(os.Stdin)
					cont.Scan()
					rs.ViewResult(result)

				} else {
					fmt.Println("not a valid file number")
				}
			}
		}
	}
}
