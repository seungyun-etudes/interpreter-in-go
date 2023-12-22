package lexer

import "monkey/token"

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
		if l.peekChar() == '=' {
			ch := l.char
			l.readChar()
			literal := string(ch) + string(l.char)
			searched = token.New(token.EQUAL, literal)
		} else {
			searched = token.New(token.ASSIGN, string(l.char))
		}
	case ';':
		searched = token.New(token.SEMICOLON, string(l.char))
	case '(':
		searched = token.New(token.LPAREN, string(l.char))
	case ')':
		searched = token.New(token.RPAREN, string(l.char))
	case '{':
		searched = token.New(token.LBRACE, string(l.char))
	case '}':
		searched = token.New(token.RBRACE, string(l.char))
	case ',':
		searched = token.New(token.COMMA, string(l.char))
	case '+':
		searched = token.New(token.PLUS, string(l.char))
	case '-':
		searched = token.New(token.MINUS, string(l.char))
	case '*':
		searched = token.New(token.ASTERISK, string(l.char))
	case '/':
		searched = token.New(token.SLASH, string(l.char))
	case '!':
		if l.peekChar() == '=' {
			ch := l.char
			l.readChar()
			literal := string(ch) + string(l.char)
			searched = token.New(token.NOT_EQUAL, literal)
		} else {
			searched = token.New(token.BANG, string(l.char))
		}
	case '<':
		searched = token.New(token.LESS, string(l.char))
	case '>':
		searched = token.New(token.GREATER, string(l.char))
	case '"':
		searched = token.New(token.STRING, l.readString())
	case 0:
		searched = token.New(token.EOF, "")
	default:
		if isLetter(l.char) {

			s := l.readIdentifier()
			return token.New(token.FindStringType(s), s)
		} else if isDigit(l.char) {
			return token.New(token.NUMBER, l.readNumber())
		} else {
			searched = token.New(token.ILLEGAL, string(l.char))
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
	start := l.current + 1
	for {
		l.readChar()
		if l.char == '"' || l.char == 0 {
			break
		}
	}
	return l.input[start:l.current]
}

func (l *Lexer) readIdentifier() string {
	start := l.current
	for isLetter(l.char) || isDigit(l.char) {
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

func (l *Lexer) peekChar() byte {
	if l.peek >= len(l.input) {
		return 0
	}

	return l.input[l.peek]
}

func isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '_'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
