package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func rs() {
	args := []string{"-t", "-b", "dir:/home/schulle/ownCloud/rezepte"}
	fmt.Print("Enter search: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	search := scanner.Text()

	args = append(args, search)
	cmd := exec.Command("recoll", args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
}

func main() {
	rs()
}
