package main

import (
	"fmt"
	"monkey/repl"
	"os"
)

func main() {
	fmt.Printf("Monkey Pogramming Language Interpreter \n\n")
	repl.Start(os.Stdin, os.Stdout)
}
