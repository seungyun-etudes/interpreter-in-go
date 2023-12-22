package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	ILLEGAL   = "ILLEGAL"
	EOF       = "EOF"
	ID        = "ID"
	NUMBER    = "NUMBER"
	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	ASTERISK  = "*"
	SLASH     = "/"
	BANG      = "!"
	LESS      = "<"
	GREATER   = ">"
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	EQUAL     = "=="
	NOT_EQUAL = "!="
	FUNCTION  = "FUNCTION"
	LET       = "LET"
	RETURN    = "RETURN"
	TRUE      = "TRUE"
	FALSE     = "FALSE"
	IF        = "IF"
	ELSE      = "ELSE"
	STRING    = "STRING"
)

var keywords = map[string]Type{
	"function": FUNCTION,
	"let":      LET,
	"return":   RETURN,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
}

func New(tokenType Type, literal string) Token {
	return Token{Type: tokenType, Literal: literal}
}

func FindStringType(s string) Type {
	if t, ok := keywords[s]; ok {
		return t
	}
	return ID
}
