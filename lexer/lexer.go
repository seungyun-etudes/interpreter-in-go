package lexer

import token "monkey/token"

type Lexer struct {
	input   string
	current int
	peek    int
	char    byte
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readChar()
	return lexer
}

func (l *Lexer) NextToken() token.Token {
	var searched token.Token

	l.skipWhitespace()

	switch l.char {
	case '=':
		searched = token.Token{Type: token.ASSIGN, Literal: string(l.char)}
	case ';':
		searched = token.Token{Type: token.SEMICOLON, Literal: string(l.char)}
	case '(':
		searched = token.Token{Type: token.LPAREN, Literal: string(l.char)}
	case ')':
		searched = token.Token{Type: token.RPAREN, Literal: string(l.char)}
	case '{':
		searched = token.Token{Type: token.LBRACE, Literal: string(l.char)}
	case '}':
		searched = token.Token{Type: token.RBRACE, Literal: string(l.char)}
	case ',':
		searched = token.Token{Type: token.COMMA, Literal: string(l.char)}
	case '+':
		searched = token.Token{Type: token.PLUS, Literal: string(l.char)}
	case '-':
		searched = token.Token{Type: token.MINUS, Literal: string(l.char)}
	case '*':
		searched = token.Token{Type: token.ASTERISK, Literal: string(l.char)}
	case '/':
		searched = token.Token{Type: token.SLASH, Literal: string(l.char)}
	case '!':
		searched = token.Token{Type: token.BANG, Literal: string(l.char)}
	case '<':
		searched = token.Token{Type: token.LESS, Literal: string(l.char)}
	case '>':
		searched = token.Token{Type: token.GREATER, Literal: string(l.char)}
	case 0:
		searched = token.Token{Type: token.EOF, Literal: ""}
	default:
		if isLetter(l.char) {
			s := l.readString()
			return token.Token{Type: token.FindStringType(s), Literal: s}
		} else if isDigit(l.char) {
			return token.Token{Type: token.NUMBER, Literal: l.readNumber()}
		} else {
			searched = token.Token{Type: token.ILLEGAL, Literal: string(l.char)}
		}
	}

	l.readChar()

	return searched
}

func (l *Lexer) readChar() {
	if l.peek >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.peek]
	}
	l.current = l.peek
	l.peek++
}

func (l *Lexer) readString() string {
	start := l.current
	for isLetter(l.char) {
		l.readChar()
	}
	return l.input[start:l.current]
}

func (l *Lexer) readNumber() string {
	start := l.current
	for isDigit(l.char) {
		l.readChar()
	}
	return l.input[start:l.current]
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '_'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
