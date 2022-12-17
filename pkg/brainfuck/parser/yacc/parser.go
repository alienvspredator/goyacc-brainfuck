// Package yacc is a brainfuck parser implemented with goyacc.
//
//go:generate goyacc -o yacc.go -p "brainfuck" yacc.y
package yacc

import (
	"github.com/alienvspredator/brainfuck-yacc/pkg/brainfuck/ast"
)

// Parse parses brainfuck code and returns a program AST.
func Parse(src []byte) *ast.Program {
	p := brainfuckNewParser()
	lex := &lexer{src: src}
	p.Parse(lex)

	return lex.result
}
