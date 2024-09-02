package main

import "fmt"

type LoxFunction struct {
	declaration *Function[any]
	closure     *Environment
}

func (f *LoxFunction) arity() int {
	return len(f.declaration.params)
}

func (f *LoxFunction) call(v *interpreter, arguments []any) (r any, err error) {
	env := newEnvironment(f.closure)
	for i, p := range f.declaration.params {
		err = env.define(p, arguments[i])
		if err != nil {
			return nil, err
		}
	}

	err = v.executeBlock(f.declaration.body, env)
	if err != nil {
		re, ok := err.(*ReturnError)
		if ok {
			return re.value, nil
		}
	}

	return nil, err
}

func (f *LoxFunction) bind(instance *LoxInstance) (*LoxFunction, error) {
	env := newEnvironment(f.closure)
	err := env.define(&token{tokenType: THIS, lexeme: "this"}, instance)
	if err != nil {
		return nil, err
	}

	return &LoxFunction{declaration: f.declaration, closure: env}, nil
}

func (f LoxFunction) String() string {
	return fmt.Sprintf("<fn %v>", f.declaration.name.lexeme)
}
