package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 5;
				let y = 10;
				let foobar = 838383;
				`

	newLexer := lexer.NewLexer(input)
	parser := NewParser(newLexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser, "letState")

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

func TestReturnStatements(t *testing.T) {
	input := `return 5; return 10; return 993322;`

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser, "return")

	if len(program.Statemens) != 3 {
		t.Fatalf("program.Statements does not conatin 3 statements. got=%d", len(program.Statemens))
	}

	for _, statement := range program.Statemens {
		returnStatement, ok := statement.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("statement not *ast.ReturnStatement. got=%T", statement)
			continue
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatement.TokenLiteral not 'return', got %q", returnStatement.TokenLiteral())
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

func checkParserErrors(t *testing.T, parser *Parser, from string) {
	errors := parser.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q -- %s", msg, from)
	}

	t.FailNow()
}
