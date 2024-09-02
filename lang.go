package main

import "fmt"

const (
	// Single-character tokens.
	LEFT_PAREN  = iota
	RIGHT_PAREN = iota
	LEFT_BRACE  = iota
	RIGHT_BRACE = iota
	COMMA       = iota
	DOT         = iota
	MINUS       = iota
	PLUS        = iota
	SEMICOLON   = iota
	SLASH       = iota
	STAR        = iota

	// One or two character tokens.
	BANG          = iota
	BANG_EQUAL    = iota
	EQUAL         = iota
	EQUAL_EQUAL   = iota
	GREATER       = iota
	GREATER_EQUAL = iota
	LESS          = iota
	LESS_EQUAL    = iota

	// Literals.
	IDENTIFIER = iota
	STRING     = iota
	NUMBER     = iota

	// Keywords.
	AND    = iota
	CLASS  = iota
	ELSE   = iota
	FALSE  = iota
	FUN    = iota
	FOR    = iota
	IF     = iota
	NIL    = iota
	OR     = iota
	PRINT  = iota
	RETURN = iota
	SUPER  = iota
	THIS   = iota
	TRUE   = iota
	VAR    = iota
	WHILE  = iota

	EOF = iota
)

var keywords = map[string]int{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type token struct {
	tokenType int
	lexeme    string
	literal   any
	line      int
}

func newToken(tokenType int, lexeme string, literal any, line int) *token {
	t := &token{
		tokenType: tokenType,
		lexeme:    lexeme,
		literal:   literal,
		line:      line,
	}
	return t
}

func (t *token) String() string {
	return fmt.Sprintf("%v %v %v", t.tokenType, t.lexeme, t.literal)
}

type LoxCallable interface {
	arity() int
	call(v *interpreter, arguments []any) (any, error)
}
