package parser

import (
	"fmt"
	"strconv"

	"github.com/Devansh3712/interpreter/ast"
	"github.com/Devansh3712/interpreter/lexer"
	"github.com/Devansh3712/interpreter/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	// > or <
	LESSGREATER
	SUM
	PRODUCT
	// -x or !x
	PREFIX
	// function(x)
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.GT:       LESSGREATER,
	token.LT:       LESSGREATER,
	token.GTE:      LESSGREATER,
	token.LTE:      LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.LPAREN:   CALL,
}

type (
	prefixParseFn func() ast.Expression
	// Function argument represents the left side of the
	// infix operator being parsed
	infixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	// Points to the token after the current token
	peekToken token.Token

	errors []string

	// Check if the appropriate map has a parsing function
	// associated with currToken.Type
	// All parsing functions do not advance tokens, they're
	// only associated with currToken
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	// defer untrace(trace("parseIntegerLiteral"))
	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		message := fmt.Sprintf("could not parse %q as integer", p.currToken.Literal)
		p.errors = append(p.errors, message)
		return nil
	}
	return &ast.IntegerLiteral{Token: p.currToken, Value: value}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	// defer untrace(trace("parsePrefixExpression"))
	expression := &ast.PrefixExpression{
		Token: p.currToken, Operator: p.currToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// defer untrace(trace("parseInfixExpression"))
	expression := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}

	precedence := p.currPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currToken, Value: p.currTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expression := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.currToken, Statements: []ast.Statement{},
	}
	p.nextToken()

	for !p.currTokenIs(token.RBRACE) && !p.currTokenIs(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	identifier := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	identifiers = append(identifiers, identifier)

	for p.peekTokenIs(token.COMMA) {
		// Skip the comma
		p.nextToken()
		p.nextToken()
		identifier := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		identifiers = append(identifiers, identifier)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: p.currToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	literal.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	literal.Body = p.parseBlockStatement()
	return literal
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: p.currToken, Function: function}
	expression.Arguments = p.parseCallArguments()
	return expression
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         []string{},
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}

	prefixFns := map[token.TokenType]prefixParseFn{
		token.IDENT:    p.parseIdentifier,
		token.INT:      p.parseIntegerLiteral,
		token.BANG:     p.parsePrefixExpression,
		token.MINUS:    p.parsePrefixExpression,
		token.TRUE:     p.parseBoolean,
		token.FALSE:    p.parseBoolean,
		token.LPAREN:   p.parseGroupedExpression,
		token.IF:       p.parseIfExpression,
		token.FUNCTION: p.parseFunctionLiteral,
	}
	for token, fn := range prefixFns {
		p.registerPrefix(token, fn)
	}

	infixFns := []token.TokenType{
		token.PLUS,
		token.MINUS,
		token.ASTERISK,
		token.SLASH,
		token.EQ,
		token.NEQ,
		token.LT,
		token.GT,
		token.LTE,
		token.GTE,
	}
	for _, token := range infixFns {
		p.registerInfix(token, p.parseInfixExpression)
	}
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	// nextToken method is called twice in order to set
	// both currToken and peekToken as if it run once
	// only peekToken is set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string { return p.errors }

func (p *Parser) peekError(t token.TokenType) {
	message := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, message)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	message := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, message)
}

func (p *Parser) currTokenIs(t token.TokenType) bool { return p.currToken.Type == t }
func (p *Parser) peekTokenIs(t token.TokenType) bool { return p.peekToken.Type == t }

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: p.currToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO:
	// Skipping expressions until we encounter a semicolon
	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()
	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// defer untrace(trace("parseExpression"))
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}
	leftExpression := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExpression
		}
		p.nextToken()
		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// defer untrace(trace("parseExpressionStatement"))
	statement := &ast.ExpressionStatement{Token: p.currToken}
	statement.Expression = p.parseExpression(LOWEST)

	// Semicolon is optional (for REPL)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currTokenIs(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}

	return program
}
