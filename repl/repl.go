package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		input := scanner.Scan()

		if !input {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)

		for t := l.NextToken(); t.Type != token.EOF; t = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", t)
		}
	}
}
