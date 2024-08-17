package main

import (
	"container/list"
)

type resolver struct {
	i      *interpreter
	scopes *list.List
}

func newResolver(i *interpreter) *resolver {
	r := &resolver{i, list.New()}

	return r
}

func (r *resolver) visitAssignExpr(e *Assign[any]) (any, error) {
	result, err := r.resolveExpression(e.value)
	if err != nil {
		return nil, err
	}

	r.resolveLocal(e, e.name)

	return result, nil
}

func (r *resolver) visitBinaryExpr(e *Binary[any]) (any, error) {
	_, err := r.resolveExpression(e.left)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveExpression(e.right)
	return nil, err
}

func (r *resolver) visitBlockStmt(stmt *Block[any]) error {
	r.beginScope()
	err := r.resolve(stmt.statements)
	if err != nil {
		return err
	}
	r.endScope()
	return nil
}

func (r *resolver) visitCallExpr(e *Call[any]) (any, error) {
	_, err := r.resolveExpression(e.callee)
	if err != nil {
		return nil, err
	}

	for _, argument := range e.arguments {
		_, err = r.resolveExpression(argument)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *resolver) visitExpressionStmt(stmt *Expression[any]) error {
	_, err := r.resolveExpression(stmt.expression)
	return err
}

func (r *resolver) visitFunctionStmt(funcStmt *Function[any]) error {
	r.declare(funcStmt.name)
	r.define(funcStmt.name)

	err := r.resolveFunction(funcStmt)
	return err
}

func (r *resolver) visitGroupingExpr(e *Grouping[any]) (any, error) {
	_, err := r.resolveExpression(e.expression)
	return nil, err
}

func (r *resolver) visitIfStmt(stmt *If[any]) error {
	_, err := r.resolveExpression(stmt.condition)
	if err != nil {
		return err
	}

	err = r.resolveStatement(stmt.thenBranch)
	if err != nil {
		return err
	}

	if stmt.elseBranch != nil {
		err = r.resolveStatement(stmt.elseBranch)
	}
	return err
}

func (r *resolver) visitLiteralExpr(e *Literal[any]) (any, error) {
	return nil, nil
}

func (r *resolver) visitLogicalExpr(e *Logical[any]) (any, error) {
	_, err := r.resolveExpression(e.left)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveExpression(e.right)
	return nil, err
}

func (r *resolver) visitPrintStmt(stmt *Print[any]) error {
	_, err := r.resolveExpression(stmt.expression)
	return err
}

func (r *resolver) visitReturnStmt(stmt *Return[any]) (err error) {
	if stmt.value != nil {
		_, err = r.resolveExpression(stmt.value)
	}

	return err
}

func (r *resolver) visitUnaryExpr(e *Unary[any]) (any, error) {
	_, err := r.resolveExpression(e.right)
	return nil, err
}

func (r *resolver) visitVarStmt(s *Var[any]) (err error) {
	r.declare(s.name)
	if s.initializer != nil {
		_, err = r.resolveExpression(s.initializer)
		if err != nil {
			return
		}
	}
	r.define(s.name)

	return
}

func (r *resolver) visitVariableExpr(e *Variable[any]) (any, error) {
	if r.scopes.Len() > 0 {
		scope := r.scopes.Back().Value.(map[string]bool)
		v, ok := scope[e.name.lexeme]
		if ok && !v {
			return nil, &ParseError{e.name, "Can't read local variable in its own initializer."}
		}
	}

	r.resolveLocal(e, e.name)
	return nil, nil
}

func (r *resolver) visitWhileStmt(whileStmt *While[any]) (err error) {
	_, err = r.resolveExpression(whileStmt.condition)
	if err != nil {
		return
	}

	err = r.resolveStatement(whileStmt.body)
	return
}

func (r *resolver) resolve(statements []Stmt[any]) (err error) {
	for _, s := range statements {
		err = r.resolveStatement(s)
		if err != nil {
			return
		}
	}

	return nil
}

func (r *resolver) resolveStatement(stmt Stmt[any]) error {
	return stmt.accept(r)
}

func (r *resolver) resolveExpression(expr Expr[any]) (any, error) {
	return expr.accept(r)
}

func (r *resolver) resolveLocal(expr Expr[any], name *token) {
	scopesSize := r.scopes.Len()

	i := scopesSize - 1
	for elem := r.scopes.Back(); elem != nil; elem = elem.Prev() {
		scope := elem.Value.(map[string]bool)

		_, ok := scope[name.lexeme]
		if ok {
			depth := scopesSize - 1 - i
			r.i.resolve(expr, depth)
			return
		}
		i = i - 1
	}
}

func (r *resolver) resolveFunction(function *Function[any]) error {
	r.beginScope()
	for _, param := range function.params {
		r.declare(param)
		r.define(param)
	}
	err := r.resolve(function.body)
	if err != nil {
		return err
	}
	r.endScope()

	return nil
}

func (r *resolver) beginScope() {
	r.scopes.PushBack(map[string]bool{})
}

func (r *resolver) endScope() {
	r.scopes.Remove(r.scopes.Back())
}

func (r *resolver) declare(name *token) {
	if r.scopes.Len() == 0 {
		return
	}

	scope := r.scopes.Back().Value.(map[string]bool)
	scope[name.lexeme] = false
}

func (r *resolver) define(name *token) {
	if r.scopes.Len() == 0 {
		return
	}

	scope := r.scopes.Back().Value.(map[string]bool)
	scope[name.lexeme] = true
}
