// Data structure used for internal representation of
// source code is called abstract syntax tree (AST)
package ast

import "github.com/Devansh3712/interpreter/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	// Dummy method
	statementNode()
}

// Expressions produce values
type Expression interface {
	Node
	// Dummy method
	expressionNode()
}

// Root node of every AST
// Every valid program is a series of statements
type Program struct {
	Statements []Statement
}

type Identifier struct {
	// token.IDENT
	Token token.Token
	Value string
}

// Identifier implements the Expression interface as identifiers
// in other parts do produce values
//
// let x = valueProducingIdentifier;
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// let <statement> = <expression>;
type LetStatement struct {
	// token.LET
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
