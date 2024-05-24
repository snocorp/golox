package main

import (
	"fmt"
	"strings"
)

type astPrinter struct {
	indent int
}

func (p *astPrinter) visitAssignExpr(e *Assign[string]) (string, error) {
	value, err := p.parenthesize(fmt.Sprintf("%v =", e.name.lexeme), e.value)
	if err != nil {
		return "", err
	}

	return value, nil
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

func (p *astPrinter) visitLogicalExpr(e *Logical[string]) (string, error) {
	return p.parenthesize(e.operator.lexeme, e.left, e.right)
}

func (p *astPrinter) visitUnaryExpr(e *Unary[string]) (string, error) {
	return p.parenthesize(e.operator.lexeme, e.right)
}

func (p *astPrinter) visitVariableExpr(e *Variable[string]) (string, error) {
	return e.name.lexeme, nil
}

func (p *astPrinter) visitBlockStmt(s *Block[string]) error {
	p.println("{")
	p.indent = p.indent + 2
	for _, s := range s.statements {
		err := s.accept(p)
		if err != nil {
			return err
		}
	}
	p.indent = p.indent - 2
	p.println("}")
	return nil
}

func (p *astPrinter) visitExpressionStmt(s *Expression[string]) error {
	result, err := p.parenthesize("stmt", s.expression)
	if err != nil {
		return err
	}

	p.println(result)

	return nil
}

func (p *astPrinter) visitIfStmt(ifStmt *If[string]) error {
	result, err := p.parenthesize("if", ifStmt.condition)
	if err != nil {
		return err
	}

	p.println(result)
	p.indent = p.indent + 2
	ifStmt.thenBranch.accept(p)
	p.indent = p.indent - 2

	if ifStmt.elseBranch != nil {
		result, err = p.parenthesize("else")
		if err != nil {
			return err
		}
		p.println(result)
		ifStmt.elseBranch.accept(p)
		p.indent = p.indent - 2
	}

	return nil
}

func (p *astPrinter) visitPrintStmt(s *Print[string]) error {
	result, err := p.parenthesize("print", s.expression)
	if err != nil {
		return err
	}

	p.println(result)

	return nil
}

func (p *astPrinter) visitVarStmt(s *Var[string]) error {
	result, err := p.parenthesize(fmt.Sprintf("var %v", s.name.literal), s.initializer)
	if err != nil {
		return err
	}

	p.println(result)

	return nil
}

func (p *astPrinter) visitWhileStmt(whileStmt *While[string]) error {
	result, err := p.parenthesize("while", whileStmt.condition)
	if err != nil {
		return err
	}

	p.println(result)
	p.indent = p.indent + 2
	whileStmt.body.accept(p)
	p.indent = p.indent - 2

	return nil
}

func (p *astPrinter) visitFunctionStmt(funcStmt *Function[string]) error {
	params := make([]string, len(funcStmt.params))
	for i, p := range funcStmt.params {
		params[i] = p.lexeme
	}

	result := fmt.Sprintf("%v(%v)", funcStmt.name.lexeme, strings.Join(params, ", "))

	p.println(result)
	p.indent = p.indent + 2
	funcStmt.body.accept(p)
	p.indent = p.indent - 2

	return nil
}

func (p *astPrinter) visitCallExpr(e *Call[string]) (string, error) {
	callee, err := e.callee.accept(p)
	if err != nil {
		return "", err
	}

	return p.parenthesize(callee, e.arguments...)
}

func (p *astPrinter) println(s string) {
	fmt.Printf("%v%s\n", strings.Repeat(" ", p.indent), s)
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
