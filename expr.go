package main

type Expr[T any] interface {
	accept(v Visitor[T]) (T, error)
}

type Assign[T any] struct {
	name *token
	value Expr[T]
}

func (e *Assign[T]) accept(v Visitor[T]) (T, error) {
	return v.visitAssignExpr(e)
}

type Binary[T any] struct {
	left Expr[T]
	operator *token
	right Expr[T]
}

func (e *Binary[T]) accept(v Visitor[T]) (T, error) {
	return v.visitBinaryExpr(e)
}

type Call[T any] struct {
	callee Expr[T]
	paren *token
	arguments []Expr[T]
}

func (e *Call[T]) accept(v Visitor[T]) (T, error) {
	return v.visitCallExpr(e)
}

type Get[T any] struct {
	object Expr[T]
	name *token
}

func (e *Get[T]) accept(v Visitor[T]) (T, error) {
	return v.visitGetExpr(e)
}

type Grouping[T any] struct {
	expression Expr[T]
}

func (e *Grouping[T]) accept(v Visitor[T]) (T, error) {
	return v.visitGroupingExpr(e)
}

type Literal[T any] struct {
	value any
}

func (e *Literal[T]) accept(v Visitor[T]) (T, error) {
	return v.visitLiteralExpr(e)
}

type Logical[T any] struct {
	left Expr[T]
	operator *token
	right Expr[T]
}

func (e *Logical[T]) accept(v Visitor[T]) (T, error) {
	return v.visitLogicalExpr(e)
}

type Set[T any] struct {
	object Expr[T]
	name *token
	value Expr[T]
}

func (e *Set[T]) accept(v Visitor[T]) (T, error) {
	return v.visitSetExpr(e)
}

type Super[T any] struct {
	keyword *token
	method *token
}

func (e *Super[T]) accept(v Visitor[T]) (T, error) {
	return v.visitSuperExpr(e)
}

type This[T any] struct {
	keyword *token
}

func (e *This[T]) accept(v Visitor[T]) (T, error) {
	return v.visitThisExpr(e)
}

type Unary[T any] struct {
	operator *token
	right Expr[T]
}

func (e *Unary[T]) accept(v Visitor[T]) (T, error) {
	return v.visitUnaryExpr(e)
}

type Variable[T any] struct {
	name *token
}

func (e *Variable[T]) accept(v Visitor[T]) (T, error) {
	return v.visitVariableExpr(e)
}

