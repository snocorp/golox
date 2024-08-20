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
	class *LoxClass
}

func newLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{class: class}
}

func (i LoxInstance) String() string {
	return fmt.Sprintf("<instance %v>", i.class.name)
}
