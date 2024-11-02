// Lexical analysis is the transformation from source code
// to tokens. Tokens itself are small, easily categorizable
// data structures that are fed to the parser, which does
// the second transformation and turns tokens into an
// abstract syntax tree (AST)
package lexer

import "github.com/Devansh3712/interpreter/token"

type Lexer struct {
	input string
	// Current position in input (points to current char)
	position int
	// Current reading position in input (after current char)
	readPosition int
	// Current char under examination
	char byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// Peek the next character without incrementing
// the read position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\r' || l.char == '\n' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}

func isLetter(char byte) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || char == '_'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.char) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.char) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.char == '"' || l.char == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.char {
	case '=':
		if l.peekChar() == '=' {
			char := l.char
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(char) + string(l.char)}
		} else {
			tok = newToken(token.ASSIGN, l.char)
		}
	case '+':
		tok = newToken(token.PLUS, l.char)
	case '-':
		tok = newToken(token.MINUS, l.char)
	case '!':
		if l.peekChar() == '=' {
			char := l.char
			l.readChar()
			tok = token.Token{Type: token.NEQ, Literal: string(char) + string(l.char)}
		} else {
			tok = newToken(token.BANG, l.char)
		}
	case '/':
		tok = newToken(token.SLASH, l.char)
	case '*':
		tok = newToken(token.ASTERISK, l.char)
	case '<':
		if l.peekChar() == '=' {
			char := l.char
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(char) + string(l.char)}
		} else {
			tok = newToken(token.LT, l.char)
		}
	case '>':
		if l.peekChar() == '=' {
			char := l.char
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(char) + string(l.char)}
		} else {
			tok = newToken(token.GT, l.char)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case ',':
		tok = newToken(token.COMMA, l.char)
	case '(':
		tok = newToken(token.LPAREN, l.char)
	case ')':
		tok = newToken(token.RPAREN, l.char)
	case '{':
		tok = newToken(token.LBRACE, l.char)
	case '}':
		tok = newToken(token.RBRACE, l.char)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.char) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	}

	l.readChar()
	return tok
}
