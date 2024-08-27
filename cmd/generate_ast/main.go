package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: generate_ast <output directory>")
		os.Exit(64)
	}
	outputDir := args[0]

	err := defineAst(outputDir, "Expr", "(T, error)", []string{
		"Assign   : name *token, value Expr[T]",
		"Binary   : left Expr[T], operator *token, right Expr[T]",
		"Call     : callee Expr[T], paren *token, arguments []Expr[T]",
		"Get      : object Expr[T], name *token",
		"Grouping : expression Expr[T]",
		"Literal  : value any",
		"Logical  : left Expr[T], operator *token, right Expr[T]",
		"Unary    : operator *token, right Expr[T]",
		"Variable : name *token",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(65)
	}

	err = defineAst(outputDir, "Stmt", "error", []string{
		"Block      : statements []Stmt[T]",
		"Class      : name *token, methods []*Function[T]",
		"Expression : expression Expr[T]",
		"Function   : name *token, params []*token, body []Stmt[T]",
		"If         : condition Expr[T], thenBranch Stmt[T], elseBranch Stmt[T]",
		"Print      : expression Expr[T]",
		"Return     : keyword *token, value Expr[T]",
		"Var        : name *token, initializer Expr[T]",
		"While      : condition Expr[T], body Stmt[T]",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(65)
	}
}

func defineAst(outputDir, baseName, acceptReturn string, types []string) error {
	fileName := strings.ToLower(baseName) + ".go"
	filePath := path.Join(outputDir, fileName)

	f, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	lines := []string{
		"package main",
		"",
		fmt.Sprintf("type %s[T any] interface {", baseName),
		fmt.Sprintf("\taccept(v Visitor[T]) %s", acceptReturn),
		"}",
		"",
	}

	for _, exprType := range types {
		parts := strings.Split(exprType, ":")
		exprName := strings.TrimSpace(parts[0])
		lines = append(lines, fmt.Sprintf("type %s[T any] struct {", exprName))

		parts = strings.Split(strings.TrimSpace(parts[1]), ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			lines = append(lines, "\t"+part)
		}

		lines = append(lines,
			"}",
			"",
			fmt.Sprintf("func (e *%s[T]) accept(v Visitor[T]) %s {", exprName, acceptReturn),
			fmt.Sprintf("\treturn v.visit%s%s(e)", exprName, baseName),
			"}",
			"",
		)
	}

	w := bufio.NewWriter(f)
	for _, line := range lines {
		_, err = w.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

/**
package main

type Expr interface {
	accept(v Visitor)
}

type Binary struct {
	left     Expr
	operator token
	right    Expr
}

func (e *Binary) accept(v Visitor) {
	v.visitBinaryExpr(e)
}
*/
