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
		fmt.Println("Enter a key:")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		key := input.Text()
		fmt.Println(key)
		fmt.Printf("%T\n", key)
		// if err != nil {
		// 	fmt.Fprintln(os.Stderr, "Error getting key")
		// 	os.Exit(1)
		// }

		switch key {
		case "q":
		 	os.Exit(0)
		 case "":
		 	fmt.Println("pressed ENTER")
			// continue
		default:
			// 	print file
			i, err := strconv.Atoi(key)
			fmt.Println("i:", i)
			if err == nil {
				// fmt.Println(i)
				if i >= 1 && i <= len(result) {
					fmt.Println(result[i-1])
					rs.FileConcat(result, i-1)
				} else {
					fmt.Println("not a valid index")
					// os.Exit(1)
				}
			}
		}
	}
}
