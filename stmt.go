package main

type Stmt[T any] interface {
	accept(v Visitor[T]) error
}

type Expression[T any] struct {
	expression Expr[T]
}

func (e *Expression[T]) accept(v Visitor[T]) error {
	return v.visitExpressionStmt(e)
}

type Print[T any] struct {
	expression Expr[T]
}

func (e *Print[T]) accept(v Visitor[T]) error {
	return v.visitPrintStmt(e)
}

type Var[T any] struct {
	name *token
	initializer Expr[T]
}

func (e *Var[T]) accept(v Visitor[T]) error {
	return v.visitVarStmt(e)
}

