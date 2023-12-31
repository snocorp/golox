package main

import (
	"fmt"
)

type RuntimeError struct {
	t       *token
	message string
}

func (err *RuntimeError) Error() string {
	return fmt.Sprintf("[line %v] Runtime Error: %s", err.t.line, err.message)
}

type interpreter struct {
	env *Environment
}

func newInterpreter() *interpreter {
	return &interpreter{
		env: &Environment{values: map[string]any{}},
	}
}

func (v *interpreter) interpret(statements []Stmt[any]) {
	for _, s := range statements {
		err := v.execute(s)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (v *interpreter) visitAssignExpr(e *Assign[any]) (any, error) {
	value, err := v.evaluate(e.value)
	if err != nil {
		return nil, err
	}

	err = v.env.assign(e.name, value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (v *interpreter) visitBinaryExpr(e *Binary[any]) (any, error) {
	left, err := v.evaluate(e.left)
	if err != nil {
		return nil, err
	}
	right, err := v.evaluate(e.right)
	if err != nil {
		return nil, err
	}

	switch e.operator.tokenType {
	case GREATER, GREATER_EQUAL, LESS, LESS_EQUAL, MINUS, SLASH, STAR:
		leftValue, leftOk := left.(float64)
		rightValue, rightOk := right.(float64)
		if !leftOk || !rightOk {
			return nil, &RuntimeError{t: e.operator, message: "Operands must be numbers."}
		}

		switch e.operator.tokenType {
		case GREATER:
			return leftValue > rightValue, nil
		case GREATER_EQUAL:
			return leftValue >= rightValue, nil
		case LESS:
			return leftValue < rightValue, nil
		case LESS_EQUAL:
			return leftValue <= rightValue, nil
		case MINUS:
			return leftValue - rightValue, nil
		case SLASH:
			return leftValue / rightValue, nil
		case STAR:
			return leftValue * rightValue, nil
		}
	case PLUS:
		leftValue, leftOk := left.(float64)
		rightValue, rightOk := right.(float64)
		if leftOk && rightOk {
			return leftValue + rightValue, nil
		}

		leftString, leftOk := left.(string)
		rightString, rightOk := right.(string)
		if leftOk && rightOk {
			return leftString + rightString, nil
		}

		return nil, &RuntimeError{t: e.operator, message: "Operands must be two numbers or two strings."}
	case BANG_EQUAL:
		return !isEqual(left, right), nil
	case EQUAL_EQUAL:
		return isEqual(left, right), nil
	}

	// Unreachable
	return nil, &RuntimeError{t: e.operator, message: "Unexpected binary expression"}
}

func (v *interpreter) visitGroupingExpr(e *Grouping[any]) (any, error) {
	return v.evaluate(e.expression)
}

func (v *interpreter) visitLiteralExpr(e *Literal[any]) (any, error) {
	return e.value, nil
}

func (v *interpreter) visitUnaryExpr(e *Unary[any]) (any, error) {
	right, err := v.evaluate(e.right)
	if err != nil {
		return nil, err
	}

	switch e.operator.tokenType {
	case BANG:
		return !isTruthy(right), nil
	case MINUS:
		value := right.(float64)
		return -value, nil
	}

	// Unreachable.
	return nil, &RuntimeError{t: e.operator, message: "Unexpected unary expression"}
}

func (v *interpreter) visitVariableExpr(e *Variable[any]) (any, error) {
	return v.env.get(e.name)
}

func (v *interpreter) visitExpressionStmt(stmt *Expression[any]) error {
	_, err := v.evaluate(stmt.expression)

	return err
}

func (v *interpreter) visitPrintStmt(stmt *Print[any]) error {
	value, err := v.evaluate(stmt.expression)
	if err != nil {
		return err
	}

	fmt.Println(stringify(value))
	return nil
}

func (v *interpreter) visitVarStmt(s *Var[any]) (err error) {
	var value any
	if s.initializer != nil {
		value, err = v.evaluate(s.initializer)
		if err != nil {
			return err
		}
	}

	err = v.env.define(s.name, value)
	return err
}

func (v *interpreter) evaluate(e Expr[any]) (any, error) {
	return e.accept(v)
}

func (v *interpreter) execute(stmt Stmt[any]) error {
	return stmt.accept(v)
}

func stringify(object any) string {
	if object == nil {
		return "nil"
	}

	return fmt.Sprintf("%v", object)
}

func isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return a == b
}

func isTruthy(object any) bool {
	if object == nil {
		return false
	}

	boolValue, ok := object.(bool)
	if ok {
		return boolValue
	}

	return true
}
