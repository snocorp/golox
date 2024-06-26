package main

import "fmt"

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func newEnvironment(enclosing *Environment) *Environment {
	return &Environment{values: map[string]any{}, enclosing: enclosing}
}

func (e *Environment) define(name *token, value any) error {
	_, ok := e.values[name.lexeme]
	if ok {
		return &RuntimeError{t: name, message: "Variable already defined '" + name.lexeme + "'."}
	}

	e.values[name.lexeme] = value

	return nil
}

func (e *Environment) get(name *token) (any, error) {
	value, ok := e.values[name.lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.get(name)
		}

		return nil, &RuntimeError{t: name, message: "Undefined variable '" + name.lexeme + "'."}
	}

	return value, nil
}

func (e *Environment) assign(name *token, value any) error {
	_, ok := e.values[name.lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.assign(name, value)
		}

		return &RuntimeError{t: name, message: "Undefined variable '" + name.lexeme + "'."}
	}

	e.values[name.lexeme] = value

	return nil
}

func (e *Environment) String() string {
	return fmt.Sprintf("%v", e.values)
}
