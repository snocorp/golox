package main

type Environment struct {
	values map[string]any
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
		return nil, &RuntimeError{t: name, message: "Undefined variable '" + name.lexeme + "'."}
	}

	return value, nil
}

func (e *Environment) assign(name *token, value any) error {
	_, ok := e.values[name.lexeme]
	if !ok {
		return &RuntimeError{t: name, message: "Undefined variable '" + name.lexeme + "'."}
	}

	e.values[name.lexeme] = value

	return nil
}
