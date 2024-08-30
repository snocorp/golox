package main

import "fmt"

type LoxClass struct {
	name string
}

func (c *LoxClass) arity() int {
	return 0
}

func (c *LoxClass) call(v *interpreter, arguments []any) (any, error) {
	return newLoxInstance(c), nil
}

func (c LoxClass) String() string {
	return fmt.Sprintf("<class %v>", c.name)
}

type LoxInstance struct {
	class  *LoxClass
	fields map[string]any
}

func newLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{class: class, fields: map[string]any{}}
}

func (i LoxInstance) String() string {
	return fmt.Sprintf("<instance %v>", i.class.name)
}

func (i *LoxInstance) get(name *token) (any, error) {
	field, ok := i.fields[name.lexeme]
	if ok {
		return field, nil
	}

	return nil, &RuntimeError{t: name, message: fmt.Sprintf("Undefined property '%v'", name.lexeme)}
}

func (i *LoxInstance) set(name *token, value any) {
	i.fields[name.lexeme] = value
}
