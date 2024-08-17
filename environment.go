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

func (e *Environment) getAt(distance int, name string) (any, bool) {
	value, ok := e.ancestor(distance).values[name]
	return value, ok
}

func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}

	return environment
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

func (e *Environment) assignAt(distance int, name *token, value any) {
	e.ancestor(distance).values[name.lexeme] = value
}

func (e *Environment) String() string {
	return fmt.Sprintf("%v", e.values)
}
