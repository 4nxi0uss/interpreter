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

func TestIdentifierExpression(test *testing.T) {
	input := "foobar;"

	lex := lexer.NewLexer(input)
	par := NewParser(lex)
	program := par.ParseProgram()
	checkParserErrors(test, par, "Expression")

	if len(program.Statemens) != 1 {
		test.Fatalf("program has not enough statements. got=%d", len(program.Statemens))
	}

	stmt, ok := program.Statemens[0].(*ast.ExpressionStatement)

	if !ok {
		test.Fatalf("program.Statements[0] is not ast.ExptessionStatement. got=%T", program.Statemens[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		test.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		test.Fatalf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		test.Fatalf("ident.TokenLieral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	lexer := lexer.NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser, "integer literal expression")

	if len(program.Statemens) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statemens))
	}

	statem, ok := program.Statemens[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statemens[0] is not ast.ExpressionStatement. got=%s", program.Statemens[0])
	}

	literal, ok := statem.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", statem.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
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
