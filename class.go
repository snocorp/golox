package main

import "fmt"

type LoxClass struct {
	name       string
	superclass *LoxClass
	methods    map[string]*LoxFunction
}

func newLoxClass(name string, superclass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{name: name, superclass: superclass, methods: methods}
}

func (c *LoxClass) arity() int {
	initializer := c.findMethod("init")
	if initializer != nil {
		return initializer.arity()
	}
	return 0
}

func (c *LoxClass) call(v *interpreter, arguments []any) (any, error) {
	instance := newLoxInstance(c)
	initializer := c.findMethod("init")
	if initializer != nil {
		f, err := initializer.bind(instance)
		if err != nil {
			return nil, err
		}

		_, err = f.call(v, arguments)
		if err != nil {
			return nil, err
		}
	}
	return instance, nil
}

func (c LoxClass) findMethod(name string) *LoxFunction {
	m, ok := c.methods[name]
	if ok {
		return m
	}

	if c.superclass != nil {
		return c.superclass.findMethod(name)
	}

	return nil
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

	method := i.class.findMethod(name.lexeme)
	if method != nil {
		method, err := method.bind(i)
		if err != nil {
			return nil, err
		}

		return method, nil
	}

	return nil, &RuntimeError{t: name, message: fmt.Sprintf("Undefined property '%v'", name.lexeme)}
}

func (i *LoxInstance) set(name *token, value any) {
	i.fields[name.lexeme] = value
}
