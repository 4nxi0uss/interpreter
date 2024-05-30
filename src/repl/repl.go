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
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		lexerNew := lexer.NewLexer(line)

		for tok := lexerNew.NextToken(); tok.Type != token.EOF; tok = lexerNew.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
