package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	currentToken token.Token
	peekToken    token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

/*Parser constructor */
func NewParser(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer, errors: []string{}}

	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	parser.registerPrefix(token.IDENT, parser.parserIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)

	parser.nextToken()
	parser.nextToken()

	return parser
}

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) nextToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statemens = []ast.Statement{}

	for parser.currentToken.Type != token.EOF {
		statement := parser.parseStatement()

		if statement != nil {
			program.Statemens = append(program.Statemens, statement)
		}

		parser.nextToken()
	}

	return program
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.currentToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parserExpression(precedence int) ast.Expression {
	prefix := parser.prefixParseFns[parser.currentToken.Type]

	if prefix == nil {
		return nil
	}

	leftExp := prefix()

	return leftExp
}

func (parser *Parser) parserIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: parser.currentToken}

	if !parser.expectPeek(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	for !parser.currentTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: parser.currentToken}

	parser.nextToken()

	for !parser.currentTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: parser.currentToken}

	statement.Expression = parser.parserExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: parser.currentToken}

	value, err := strconv.ParseInt(parser.currentToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", parser.currentToken.Literal)
		parser.errors = append(parser.errors, msg)

		return nil
	}

	lit.Value = value

	return lit
}

func (parser *Parser) currentTokenIs(tokenType token.TokenType) bool {
	return parser.currentToken.Type == tokenType
}

func (parser *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return parser.peekToken.Type == tokenType
}

func (parser *Parser) peekError(tokenType token.TokenType) {
	message := fmt.Sprintf("expected next token to be %s, got %s instead", tokenType, parser.peekToken.Type)

	parser.errors = append(parser.errors, message)
}

func (parser *Parser) expectPeek(tokenType token.TokenType) bool {
	if parser.peekTokenIs(tokenType) {
		parser.nextToken()

		return true
	} else {
		parser.peekError(tokenType)

		return false
	}
}

func (parser *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	parser.prefixParseFns[tokenType] = fn

}

func (parser *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	parser.infixParseFns[tokenType] = fn
}
