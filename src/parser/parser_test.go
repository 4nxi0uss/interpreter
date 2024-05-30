package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `let x  5;
				let y = 10;
				let foobar = 838383;
				`

	newLexer := lexer.NewLexer(input)
	parser := NewParser(newLexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statemens) != 3 {
		t.Fatalf("program.Statemens does not contain 3 statements. got =%d", len(program.Statemens))
	}

	tests := []struct{ expectedIdentifier string }{
		{"x"}, {"y"}, {"foobar"},
	}

	for i, test := range tests {
		statement := program.Statemens[i]

		if !testLetStatement(t, statement, test.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("statement.TokenLiteral() not 'let'. got=%q", statement.TokenLiteral())

		return false
	}

	letStatement, ok := statement.(*ast.LetStatement)

	if !ok {
		t.Errorf("statement not *ast.LetStatement. got=%T", statement)

		return false
	}

	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value not '%s'. got=%s", name, letStatement.Name.TokenLiteral())

		return false
	}

	return true
}

func checkParserErrors(t *testing.T, parser *Parser) {
	errors := parser.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}
