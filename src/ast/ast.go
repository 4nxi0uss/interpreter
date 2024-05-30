package ast

import "monkey/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statemens []Statement
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Valeu Expression
}

type Identifier struct {
	Token token.Token
	Value string
}

func (letStatement *LetStatement) statementNode()       {}
func (letStatement *LetStatement) TokenLiteral() string { return letStatement.Token.Literal }

func (id *Identifier) expressionNode()      {}
func (id *Identifier) TokenLiteral() string { return id.Token.Literal }

func (program *Program) TokenLiteral() string {
	if len(program.Statemens) > 0 {
		return program.Statemens[0].TokenLiteral()
	} else {
		return ""
	}
}
