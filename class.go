package main

import "fmt"

type LoxClass struct {
	name string
}

func (c LoxClass) String() string {
	return fmt.Sprintf("<class %v>", c.name)
}
