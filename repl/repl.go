// REPL (Read Eval Print Loop) reads an input, sends it to
// the interpreter for evaluation, prints the output of the
// interpreter and starts again
package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Devansh3712/interpreter/evaluator"
	"github.com/Devansh3712/interpreter/lexer"
	"github.com/Devansh3712/interpreter/parser"
)

const PROMPT = ">> "

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, " parser errors:\n")
	for _, message := range errors {
		io.WriteString(out, "\t"+message+"\n")
	}
}

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
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
