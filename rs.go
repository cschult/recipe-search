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

func lslah() {
	cmd := exec.Command("ls", "-lah")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}

func main() {
	rs()
	// lslah()
}
