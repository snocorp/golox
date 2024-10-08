package main

import (
	"fmt"
)

type ParseError struct {
	t       *token
	message string
}

func (err *ParseError) Error() string {
	var where string
	if err.t.tokenType == EOF {
		where = " at end"
	} else {
		where = " at " + err.t.lexeme
	}
	return fmt.Sprintf("[line %v] Error%s: %v\n", err.t.line, where, err.message)
}

type Parser[T any] struct {
	tokens  []*token
	current int
	errors  []*ParseError
}

func newParser[T any](tokens []*token) *Parser[T] {
	return &Parser[T]{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser[T]) parse() ([]Stmt[T], error) {
	statements := []Stmt[T]{}
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	return statements, nil
}

func (p *Parser[T]) declaration() (Stmt[T], error) {
	if p.match(CLASS) {
		return p.classDeclaration()
	}
	if p.match(FUN) {
		return p.function("function")
	}
	if p.match(VAR) {
		return p.varDeclaration()
	}

	stmt, err := p.statement()
	if err != nil {
		p.synchronize()
		return nil, err
	}

	return stmt, nil
}

func (p *Parser[T]) statement() (Stmt[T], error) {
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(LEFT_BRACE) {
		return p.block()
	}

	return p.expressionStatement()
}

func (p *Parser[T]) returnStatement() (Stmt[T], error) {
	var err error
	var value Expr[T] = nil
	keyword := p.previous()
	if !p.check(SEMICOLON) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	p.consume(SEMICOLON, "Expect ';' after return value.")
	return &Return[T]{keyword, value}, nil
}

func (p *Parser[T]) function(kind string) (*Function[T], error) {
	name, err := p.consume(IDENTIFIER, "Expect "+kind+" name.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_PAREN, "Expect '(' after "+kind+" name.")
	if err != nil {
		return nil, err
	}

	parameters := []*token{}
	if !p.check(RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				return nil, p.error(p.peek(), "Can't have more than 255 parameters.")
			}

			paramName, err := p.consume(IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, paramName)

			if !p.match(COMMA) {
				break
			}
		}
	}
	_, err = p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_BRACE, "Expect '{' before "+kind+" body.")
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &Function[T]{
		name:   name,
		params: parameters,
		body:   body.statements,
	}, nil
}

func (p *Parser[T]) forStatement() (Stmt[T], error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer Stmt[T]
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer, err = p.varDeclaration()
	} else {
		initializer, err = p.expressionStatement()
	}
	if err != nil {
		return nil, err
	}

	var condition Expr[T]
	if !p.check(SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment Expr[T]
	if !p.check(RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = &Block[T]{
			statements: []Stmt[T]{
				body,
				&Expression[T]{increment},
			},
		}
	}

	if condition == nil {
		condition = &Literal[T]{true}
	}
	body = &While[T]{condition, body}

	if initializer != nil {
		body = &Block[T]{statements: []Stmt[T]{initializer, body}}
	}

	return body, nil
}

func (p *Parser[T]) whileStatement() (Stmt[T], error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &While[T]{condition: condition, body: body}, nil
}

func (p *Parser[T]) ifStatement() (Stmt[T], error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch Stmt[T]
	if p.match(ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &If[T]{condition, thenBranch, elseBranch}, nil
}

func (p *Parser[T]) block() (*Block[T], error) {
	statements := []Stmt[T]{}

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		s, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, s)
	}

	_, err := p.consume(RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return &Block[T]{statements: statements}, nil
}

func (p *Parser[T]) classDeclaration() (Stmt[T], error) {
	name, err := p.consume(IDENTIFIER, "Expect class name.")
	if err != nil {
		return nil, err
	}

	var superclass *Variable[T]
	if p.match(LESS) {
		p.consume(IDENTIFIER, "Expect superclass name")
		superclass = &Variable[T]{name: p.previous()}
	}

	_, err = p.consume(LEFT_BRACE, "Expect '{' before class body.")
	if err != nil {
		return nil, err
	}

	methods := []*Function[T]{}
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		fn, err := p.function("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, fn)
	}

	_, err = p.consume(RIGHT_BRACE, "Expect '}' after class body.")
	if err != nil {
		return nil, err
	}

	return &Class[T]{name: name, superclass: superclass, methods: methods}, nil
}

func (p *Parser[T]) varDeclaration() (Stmt[T], error) {
	name, err := p.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr[T]
	if p.match(EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return &Var[T]{name: name, initializer: initializer}, nil
}

func (p *Parser[T]) printStatement() (Stmt[T], error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(SEMICOLON, "Expected ';' after value.")
	if err != nil {
		return nil, err
	}

	return &Print[T]{expression: value}, nil
}

func (p *Parser[T]) expressionStatement() (Stmt[T], error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}

	return &Expression[T]{expression: expr}, nil
}

func (p *Parser[T]) expression() (Expr[T], error) {
	return p.assignment()
}

func (p *Parser[T]) assignment() (Expr[T], error) {
	var err error
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		variable, ok := expr.(*Variable[T])
		if ok {
			name := variable.name
			return &Assign[T]{name: name, value: value}, nil
		}

		get, ok := expr.(*Get[T])
		if ok {
			return &Set[T]{object: get.object, name: get.name, value: value}, nil
		}

		return nil, p.error(equals, "Invalid assignment target.")
	}

	return expr, nil
}

func (p *Parser[T]) or() (Expr[T], error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}

		expr = &Logical[T]{left: expr, operator: operator, right: right}
	}

	return expr, nil
}

func (p *Parser[T]) and() (Expr[T], error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		expr = &Logical[T]{left: expr, operator: operator, right: right}
	}

	return expr, nil
}

func (p *Parser[T]) equality() (Expr[T], error) {
	e, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		e = &Binary[T]{left: e, operator: operator, right: right}
	}

	return e, nil
}

func (p *Parser[T]) comparison() (Expr[T], error) {
	e, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		e = &Binary[T]{left: e, operator: operator, right: right}
	}

	return e, nil
}

func (p *Parser[T]) term() (Expr[T], error) {
	e, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		e = &Binary[T]{left: e, operator: operator, right: right}
	}

	return e, nil
}

func (p *Parser[T]) factor() (Expr[T], error) {
	e, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		e = &Binary[T]{left: e, operator: operator, right: right}
	}

	return e, nil
}

func (p *Parser[T]) unary() (Expr[T], error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &Unary[T]{operator, right}, nil
	}

	return p.call()
}

func (p *Parser[T]) call() (Expr[T], error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for true {
		if p.match(LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(DOT) {
			name, err := p.consume(IDENTIFIER, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = &Get[T]{expr, name}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser[T]) finishCall(callee Expr[T]) (Expr[T], error) {
	arguments := []Expr[T]{}
	if !p.check(RIGHT_PAREN) {
		for {
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}

			if len(arguments) >= 255 {
				return nil, p.error(p.peek(), "Can't have more than 255 arguments.")
			}

			arguments = append(arguments, arg)
			if !p.match(COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &Call[T]{callee, paren, arguments}, nil
}

func (p *Parser[T]) primary() (Expr[T], error) {
	if p.match(FALSE) {
		return &Literal[T]{value: false}, nil
	}
	if p.match(TRUE) {
		return &Literal[T]{value: true}, nil
	}
	if p.match(NIL) {
		return &Literal[T]{value: nil}, nil
	}

	if p.match(NUMBER, STRING) {
		return &Literal[T]{value: p.previous().literal}, nil
	}

	if p.match(SUPER) {
		keyword := p.previous()
		_, err := p.consume(DOT, "Expect '.' after 'super'.")
		if err != nil {
			return nil, err
		}

		method, err := p.consume(IDENTIFIER, "Expect superclass method name.")
		if err != nil {
			return nil, err
		}

		return &Super[T]{keyword: keyword, method: method}, nil
	}

	if p.match(THIS) {
		return &This[T]{keyword: p.previous()}, nil
	}

	if p.match(IDENTIFIER) {
		return &Variable[T]{name: p.previous()}, nil
	}

	if p.match(LEFT_PAREN) {
		e, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(RIGHT_PAREN, "Expected ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &Grouping[T]{expression: e}, nil
	}

	return nil, p.error(p.previous(), "Expected expression.")
}

func (p *Parser[T]) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().tokenType == SEMICOLON {
			return
		}

		switch p.peek().tokenType {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser[T]) consume(tokenType int, message string) (*token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return nil, p.error(p.peek(), message)
}

func (p *Parser[T]) error(t *token, message string) error {
	err := &ParseError{t: t, message: message}
	p.errors = append(p.errors, err)
	return err
}

func (p *Parser[T]) match(tokenTypes ...int) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser[T]) check(tokenType int) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == tokenType
}

func (p *Parser[T]) advance() *token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *Parser[T]) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *Parser[T]) peek() *token {
	return p.tokens[p.current]
}

func (p *Parser[T]) previous() *token {
	return p.tokens[p.current-1]
}
