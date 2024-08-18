// REPL (Read Eval Print Loop) reads an input, sends it to
// the interpreter for evaluation, prints the output of the
// interpreter and starts again
package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Devansh3712/interpreter/lexer"
	"github.com/Devansh3712/interpreter/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
