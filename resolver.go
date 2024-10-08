package main

import (
	"container/list"
	"fmt"
)

const (
	FUNC_TYPE_NONE     = iota
	FUNC_TYPE_FUNCTION = iota
	FUNC_TYPE_INIT     = iota
	FUNC_TYPE_METHOD   = iota
)

const (
	CLASS_TYPE_NONE     = iota
	CLASS_TYPE_CLASS    = iota
	CLASS_TYPE_SUBCLASS = iota
)

type ResolverError struct {
	t       *token
	message string
}

func (err *ResolverError) Error() string {
	return fmt.Sprintf("[line %v] Resolver Error: %s", err.t.line, err.message)
}

type resolver struct {
	i           *interpreter
	scopes      *list.List
	inFuncType  int
	inClassType int
}

func newResolver(i *interpreter) *resolver {
	r := &resolver{i, list.New(), FUNC_TYPE_NONE, CLASS_TYPE_NONE}

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

func (r *resolver) visitClassStmt(stmt *Class[any]) error {
	inEnclosingClassType := r.inClassType
	r.inClassType = CLASS_TYPE_CLASS

	err := r.declare(stmt.name)
	if err != nil {
		return err
	}

	r.define(stmt.name)

	if stmt.superclass != nil {
		r.inClassType = CLASS_TYPE_SUBCLASS

		if stmt.superclass.name.lexeme == stmt.name.lexeme {
			return &ResolverError{t: stmt.superclass.name, message: "A class can't inherit from itself."}
		}

		_, err = r.resolveExpression(stmt.superclass)
		if err != nil {
			return err
		}

		r.beginScope()
		scope := r.scopes.Back().Value.(map[string]bool)
		scope["super"] = true
	}

	r.beginScope()
	scope := r.scopes.Back().Value.(map[string]bool)
	scope["this"] = true

	for _, method := range stmt.methods {
		funcType := FUNC_TYPE_METHOD
		if method.name.lexeme == "init" {
			funcType = FUNC_TYPE_INIT
		}
		err := r.resolveFunction(method, funcType)
		if err != nil {
			return err
		}
	}

	r.endScope()

	if stmt.superclass != nil {
		r.endScope()
	}

	r.inClassType = inEnclosingClassType

	return nil
}

func (r *resolver) visitGetExpr(expr *Get[any]) (any, error) {
	return r.resolveExpression(expr.object)
}

func (r *resolver) visitSetExpr(expr *Set[any]) (any, error) {
	var err error

	_, err = r.resolveExpression(expr.value)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveExpression(expr.object)
	return nil, err
}

func (r *resolver) visitThisExpr(e *This[any]) (any, error) {
	if r.inClassType == CLASS_TYPE_NONE {
		return nil, &ResolverError{t: e.keyword, message: "Can't use 'this' outside of a class."}
	}
	r.resolveLocal(e, e.keyword)
	return nil, nil
}

func (r *resolver) visitSuperExpr(e *Super[any]) (any, error) {
	if r.inClassType == CLASS_TYPE_NONE {
		return nil, &ResolverError{t: e.keyword, message: "Can't use 'super' outside of a class."}
	} else if r.inClassType != CLASS_TYPE_SUBCLASS {
		return nil, &ResolverError{t: e.keyword, message: "Can't use 'super' in a class with no superclass."}
	}

	r.resolveLocal(e, e.keyword)
	return nil, nil
}

func (r *resolver) visitExpressionStmt(stmt *Expression[any]) error {
	_, err := r.resolveExpression(stmt.expression)
	return err
}

func (r *resolver) visitFunctionStmt(funcStmt *Function[any]) error {
	err := r.declare(funcStmt.name)
	if err != nil {
		return err
	}
	r.define(funcStmt.name)

	err = r.resolveFunction(funcStmt, FUNC_TYPE_FUNCTION)
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
	if r.inFuncType == FUNC_TYPE_NONE {
		return &ResolverError{t: stmt.keyword, message: "Can't return from top-level code."}
	}
	if stmt.value != nil {
		if r.inFuncType == FUNC_TYPE_INIT {
			return &ResolverError{t: stmt.keyword, message: "Can't return from an initializer."}
		}

		_, err = r.resolveExpression(stmt.value)
	}

	return err
}

func (r *resolver) visitUnaryExpr(e *Unary[any]) (any, error) {
	_, err := r.resolveExpression(e.right)
	return nil, err
}

func (r *resolver) visitVarStmt(s *Var[any]) (err error) {
	err = r.declare(s.name)
	if err != nil {
		return err
	}
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

func (r *resolver) resolveFunction(function *Function[any], funcType int) (err error) {
	inEnclosingFuncType := r.inFuncType
	r.inFuncType = funcType

	r.beginScope()
	for _, param := range function.params {
		err = r.declare(param)
		if err != nil {
			return err
		}
		r.define(param)
	}
	err = r.resolve(function.body)
	if err != nil {
		return err
	}
	r.endScope()

	r.inFuncType = inEnclosingFuncType

	return nil
}

func (r *resolver) beginScope() {
	r.scopes.PushBack(map[string]bool{})
}

func (r *resolver) endScope() {
	r.scopes.Remove(r.scopes.Back())
}

func (r *resolver) declare(name *token) error {
	if r.scopes.Len() == 0 {
		return nil
	}

	scope := r.scopes.Back().Value.(map[string]bool)
	_, ok := scope[name.lexeme]
	if ok {
		return &ResolverError{t: name, message: "Already a variable with this name in this scope."}
	}
	scope[name.lexeme] = false

	return nil
}

func (r *resolver) define(name *token) {
	if r.scopes.Len() == 0 {
		return
	}

	scope := r.scopes.Back().Value.(map[string]bool)
	scope[name.lexeme] = true
}
