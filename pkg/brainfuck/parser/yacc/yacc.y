%{

package yacc

import (
	"github.com/alienvspredator/brainfuck-yacc/pkg/brainfuck/ast"
)

%}

%union{
	nodes []ast.Node
}

%type <nodes> exprs expr loop command

%%

prog: exprs { brainfucklex.(*lexer).result = &ast.Program{Body: &ast.Body{List: $1}} }

exprs: expr exprs { $$ = append($1, $2...) }
     | expr { $$ = $1 }

expr: command { $$ = $1 }
    | loop { $$ = $1 }

loop: '[' exprs ']' { $$ = append($$, &ast.Loop{Body: &ast.Body{List: $2}}) }

command: '>' { $$ = append($$, &ast.IncPtr{}) }
       | '<' { $$ = append($$, &ast.DecPtr{}) }
       | '+' { $$ = append($$, &ast.IncByte{}) }
       | '-' { $$ = append($$, &ast.DecByte{}) }
       | '.' { $$ = append($$, &ast.OutputByte{}) }
       | ',' { $$ = append($$, &ast.InputByte{}) }
