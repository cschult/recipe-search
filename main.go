package main

import (
	"bufio"
	"devmem.de/srv/git/recipe-search/rs"
	"fmt"
	"os"
	"strconv"
)


func main() {
	result := rs.Search()

	for true {
		rs.ViewResult(result)
		// fmt.Println("| q = quit | <num> = show file <num> |")
		fmt.Println("q = quit, 1 = show file #1, ...")
		fmt.Println("Enter a key:")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		key := input.Text()

		switch key {
		case "q":
		 	os.Exit(0)
		 case "":
		 	fmt.Println("pressed ENTER")
			// continue
		default:
			// 	print file
			i, err := strconv.Atoi(key)
			// fmt.Println("i:", i)
			if err == nil {
				// fmt.Println(i)
				if i >= 1 && i <= len(result) {
					// fmt.Println(result[i-1])
					rs.FileConcat(result, i-1)
				} else {
					fmt.Println("not a valid file number")
					// os.Exit(1)
				}
			}
		}
	}
}
