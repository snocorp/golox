package main

import (
	"fmt"
	"strings"
)

type astPrinter struct{}

func (p *astPrinter) print(e Expr[string]) {
	s, err := e.accept(p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s)
}

func (p *astPrinter) visitAssignExpr(e *Assign[string]) (string, error) {
	value, err := p.parenthesize("assign", e.value)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v = %v", e.name.lexeme, value), nil
}

func (p *astPrinter) visitBinaryExpr(e *Binary[string]) (string, error) {
	return p.parenthesize(e.operator.lexeme, e.left, e.right)
}

func (p *astPrinter) visitGroupingExpr(e *Grouping[string]) (string, error) {
	return p.parenthesize("group", e.expression)
}

func (p *astPrinter) visitLiteralExpr(e *Literal[string]) (string, error) {
	if e.value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", e.value), nil
}

func (p *astPrinter) visitUnaryExpr(e *Unary[string]) (string, error) {
	return p.parenthesize(e.operator.lexeme, e.right)
}

func (p *astPrinter) visitVariableExpr(e *Variable[string]) (string, error) {
	return e.name.lexeme, nil
}

func (p *astPrinter) visitExpressionStmt(s *Expression[string]) error {
	_, err := p.parenthesize("stmt", s.expression)

	return err
}

func (p *astPrinter) visitPrintStmt(s *Print[string]) error {
	_, err := p.parenthesize("print", s.expression)

	return err
}

func (p *astPrinter) visitVarStmt(s *Var[string]) error {
	_, err := p.parenthesize(fmt.Sprintf("var %v", s.name.literal), s.initializer)

	return err
}

func (p *astPrinter) parenthesize(name string, expressions ...Expr[string]) (string, error) {
	parts := []string{
		"(",
		name,
	}

	for _, e := range expressions {
		s, err := e.accept(p)
		if err != nil {
			return "", err
		}

		parts = append(parts,
			" ",
			s,
		)
	}

	parts = append(parts, ")")

	return strings.Join(parts, ""), nil
}
