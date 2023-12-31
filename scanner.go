package main

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

type scanError struct {
	line    int
	message string
}

func (err *scanError) Error() string {
	return fmt.Sprintf("[line %v]: %v", err.line, err.message)
}

type scanner struct {
	source               string
	tokens               []*token
	start, current, line int
}

func newScanner(source string) *scanner {
	return &scanner{
		source: source,
		line:   1,
	}
}

func (s *scanner) scanTokens() ([]*token, error) {
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			return s.tokens, err
		}
	}

	s.tokens = append(s.tokens, newToken(EOF, "", nil, s.line))
	return s.tokens, nil
}

func (s *scanner) scanToken() error {
	c := s.advance()
	switch c {
	case "(":
		s.addToken(LEFT_PAREN, nil)
	case ")":
		s.addToken(RIGHT_PAREN, nil)
	case "{":
		s.addToken(LEFT_BRACE, nil)
	case "}":
		s.addToken(RIGHT_BRACE, nil)
	case ",":
		s.addToken(COMMA, nil)
	case ".":
		s.addToken(DOT, nil)
	case "-":
		s.addToken(MINUS, nil)
	case "+":
		s.addToken(PLUS, nil)
	case ";":
		s.addToken(SEMICOLON, nil)
	case "*":
		s.addToken(STAR, nil)
	case "!":
		s.addToken(s.match("=", BANG_EQUAL, BANG), nil)
	case "=":
		s.addToken(s.match("=", EQUAL_EQUAL, EQUAL), nil)
	case "<":
		s.addToken(s.match("=", LESS_EQUAL, LESS), nil)
	case ">":
		s.addToken(s.match("=", GREATER_EQUAL, GREATER), nil)
	case "/":
		if s.match("/", SLASH, -1) == SLASH {
			// A comment goes until the end of the line.
			for s.peek() != "\n" && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH, nil)
		}
	case " ", "\r", "\t":
		// Ignore whitespace.
	case "\n":
		s.line += 1
	case "\"":
		s.scanString()
	default:
		if isDigit(c) {
			s.scanNumber()
		} else if isAlpha(c) {
			s.scanIdentifier()
		} else {
			return &scanError{s.line, "Unexpected character."}
		}
	}

	return nil
}

func (s *scanner) addToken(tokenType int, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, newToken(tokenType, text, literal, s.line))
}

func (s *scanner) advance() string {
	runeValue, width := utf8.DecodeRuneInString(s.source[s.current:])
	s.current += width

	return string(runeValue)
}

func (s *scanner) peek() string {
	if s.isAtEnd() {
		return ""
	}
	runeValue, _ := utf8.DecodeRuneInString(s.source[s.current:])

	return string(runeValue)
}

func (s *scanner) peekNext() string {
	if s.current+1 >= len(s.source) {
		return ""
	}
	runeValue, _ := utf8.DecodeRuneInString(s.source[s.current+1:])

	return string(runeValue)
}

func (s *scanner) match(expected string, matchTokenType, missTokenType int) int {
	if s.isAtEnd() {
		return missTokenType
	}
	runeValue, width := utf8.DecodeRuneInString(s.source[s.current:])
	if string(runeValue) != expected {
		return missTokenType
	}

	s.current += width
	return matchTokenType
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *scanner) scanString() error {
	for s.peek() != "\"" && !s.isAtEnd() {
		if s.peek() == "\n" {
			s.line += 1
		}
		s.advance()
	}

	if s.isAtEnd() {
		return &scanError{s.line, "Unterminated string."}
	}

	// The closing ".
	s.advance()

	// Trim the surrounding quotes.
	value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, value)

	return nil
}

func (s *scanner) scanNumber() error {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == "." && isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		return &scanError{s.line, "Invalid numeric value"}
	}

	s.addToken(NUMBER, value)

	return nil
}

func (s *scanner) scanIdentifier() error {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = IDENTIFIER
	}
	s.addToken(tokenType, nil)

	return nil
}

func isDigit(c string) bool {
	if len(c) != 1 {
		return false
	}
	b := c[0]
	return b >= '0' && b <= '9'
}

func isAlpha(c string) bool {
	if len(c) != 1 {
		return false
	}
	b := c[0]
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		b == '_'
}

func isAlphaNumeric(c string) bool {
	return isAlpha(c) || isDigit(c)
}
