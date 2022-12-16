%{

package brainfuck

import (
	"fmt"
	"log"
	"unicode/utf8"
)

%}

%union{
	expr Expr
}

%type <expr> command expr loop '>' '<' '+' '-' '.'

%%

programm: expr
        {
        	NewProgram(3000, $1).Run()
        	fmt.Println()
       }

expr: command expr { $$ = NewList($1, $2) }
    | command
    | loop expr { $$ = NewList($1, $2) }
    | loop

loop: '[' expr ']' { $$ = NewLoop($2) }

command: '>' { $$ = NewIncPtr() }
       | '<' { $$ = NewDecPtr() }
       | '+' { $$ = NewAdd() }
       | '-' { $$ = NewSubtract() }
       | '.' { $$ = NewOutput() }

%%

const eof = 0

type brainfuckLex struct {
	line []byte
	peek rune
}

func (x *brainfuckLex) Lex(lval *brainfuckSymType) int {
	for {
		c := x.next()
		switch c {
		case eof:
			return eof
		case ' ', '\t', '\n', '\r':
		case '<', '>', '+', '-', '.', '[', ']':
			return int(c)
		default:
			log.Printf("unrecognized character %q", c)
		}
	}
}

// Return the next rune for the lexer.
func (x *brainfuckLex) next() rune {
	if x.peek != eof {
		r := x.peek
		x.peek = eof
		return r
	}
	if len(x.line) == 0 {
		return eof
	}
	c, size := utf8.DecodeRune(x.line)
	x.line = x.line[size:]
	if c == utf8.RuneError && size == 1 {
		log.Print("invalid utf8")
		return x.next()
	}
	return c
}

// The parser calls this method on a parse error.
func (x *brainfuckLex) Error(s string) {
	log.Printf("parse error: %s", s)
}
