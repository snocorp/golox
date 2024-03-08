package main

import "time"

type clock struct{}

func (*clock) arity() int {
	return 0
}

func (*clock) call(v *interpreter, arguments []any) (any, error) {
	return float64(time.Now().UnixMilli()), nil
}
