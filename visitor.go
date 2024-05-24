package main

type Visitor[T any] interface {
	visitAssignExpr(e *Assign[T]) (T, error)
	visitBinaryExpr(e *Binary[T]) (T, error)
	visitCallExpr(e *Call[T]) (T, error)
	visitGroupingExpr(e *Grouping[T]) (T, error)
	visitLiteralExpr(e *Literal[T]) (T, error)
	visitLogicalExpr(e *Logical[T]) (T, error)
	visitUnaryExpr(e *Unary[T]) (T, error)
	visitVariableExpr(e *Variable[T]) (T, error)

	visitBlockStmt(s *Block[T]) error
	visitExpressionStmt(s *Expression[T]) error
	visitFunctionStmt(s *Function[T]) error
	visitIfStmt(ifStmt *If[T]) error
	visitPrintStmt(s *Print[T]) error
	visitVarStmt(s *Var[T]) error
	visitWhileStmt(s *While[T]) error
}
