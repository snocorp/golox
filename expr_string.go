package main

import "fmt"

func (e *Assign[T]) String() string {
	return fmt.Sprintf("%v = %v", e.name.lexeme, e.value)
}

func (e *Binary[T]) String() string {
	return fmt.Sprintf("%v %v %v", e.left, e.operator.lexeme, e.right)
}

func (e *Call[T]) String() string {
	return fmt.Sprintf("%v(%v)", e.callee, e.arguments)
}

func (e *Grouping[T]) String() string {
	return fmt.Sprintf("{ %v }", e.expression)
}

func (e *Literal[T]) String() string {
	return fmt.Sprintf("%v", e.value)
}

func (e *Logical[T]) String() string {
	return fmt.Sprintf("%v %v %v", e.left, e.operator.lexeme, e.right)
}

func (e *Unary[T]) String() string {
	return fmt.Sprintf("%v %v", e.operator, e.right)
}

func (e *Variable[T]) String() string {
	return fmt.Sprintf("%v", e.name.lexeme)
}
