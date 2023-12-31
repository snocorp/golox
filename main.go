package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type runner struct {
	hadError bool
}

func main() {
	args := os.Args[1:]
	if len(args) > 2 {
		fmt.Println("Usage: golox [run|print] [script]")
		fmt.Println(args)
		os.Exit(64)
	} else if len(args) == 2 {
		if args[0] == "print" {
			err := printFile(args[1])
			if err != nil {
				fmt.Println(err)
				os.Exit(66)
			}
		} else {
			err := runFile(args[1])
			if err != nil {
				fmt.Println(err)
				os.Exit(66)
			}
		}
	} else {
		runPrompt()
	}
}

func printFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lox := &runner{}
	source := string(bytes)
	scanner := newScanner(source)
	tokens, err := scanner.scanTokens()
	if err != nil {
		se, ok := err.(*scanError)
		if ok {
			lox.handleError(se.line, se.message)
		} else {
			lox.handleError(-1, err.Error())
		}
		return err
	}

	parser := newParser[string](tokens)
	statements, err := parser.parse()
	if err != nil {
		return err
	}

	printer := &astPrinter{}
	for _, s := range statements {
		s.accept(printer)
	}

	return nil
}

func runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lox := &runner{}
	lox.run(string(bytes))

	// Indicate an error in the exit code.
	if lox.hadError {
		os.Exit(65)
	}

	return nil
}

func runPrompt() {
	input := bufio.NewReader(os.Stdin)

	lox := &runner{}

	for {
		fmt.Print("> ")
		line, err := input.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
		}
		lox.run(line)
		lox.hadError = false
	}
}

func (lox *runner) run(source string) {
	scanner := newScanner(source)
	tokens, err := scanner.scanTokens()
	if err != nil {
		se, ok := err.(*scanError)
		if ok {
			lox.handleError(se.line, se.message)
		} else {
			lox.handleError(-1, err.Error())
		}
		return
	}

	parser := newParser[any](tokens)
	statements, err := parser.parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	inter := newInterpreter()
	inter.interpret(statements)
}

func (lox *runner) handleError(line int, message string) {
	lox.report(line, "", message)
}

func (lox *runner) report(line int, where, message string) {
	fmt.Printf("[line %v] Error%v: %v\n", line, where, message)
	lox.hadError = true
}
