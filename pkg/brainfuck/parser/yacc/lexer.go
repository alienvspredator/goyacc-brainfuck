package yacc

import (
	"log"
	"unicode/utf8"

	"github.com/alienvspredator/brainfuck-yacc/pkg/brainfuck/ast"
)

const eof = 0

type lexer struct {
	src    []byte
	peek   rune
	result *ast.Program
}

func (x *lexer) Lex(_ *brainfuckSymType) int {
	for {
		c := x.next()
		switch c {
		case eof:
			return eof
		case ' ', '\t', '\n', '\r':
		case '<', '>', '+', '-', '.', '[', ']', ',':
			return int(c)
		}
	}
}

// Return the next rune for the lexer.
func (x *lexer) next() rune {
	if x.peek != eof {
		r := x.peek
		x.peek = eof
		return r
	}
	if len(x.src) == 0 {
		return eof
	}
	c, size := utf8.DecodeRune(x.src)
	x.src = x.src[size:]
	if c == utf8.RuneError && size == 1 {
		log.Print("invalid utf8")
		return x.next()
	}
	return c
}

// The parser calls this method on a parse error.
func (x *lexer) Error(s string) {
	log.Printf("parse error: %s", s)
}
