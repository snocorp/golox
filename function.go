package main

type LoxFunction struct {
	declaration *Function[any]
}

func (f *LoxFunction) arity() int {
	return len(f.declaration.params)
}

func (f *LoxFunction) call(v *interpreter, arguments []any) (r any, err error) {
	env := newEnvironment(v.globals)
	for i, p := range f.declaration.params {
		err = env.define(p, arguments[i])
		if err != nil {
			return nil, err
		}
	}

	err = v.executeBlock(f.declaration.body.statements, env)
	return nil, err
}
