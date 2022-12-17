package runtime

import (
	"io"

	"github.com/alienvspredator/brainfuck-yacc/pkg/brainfuck/ast"
)

type Runner struct {
	runtime *runtime
	prog    *ast.Program
}

func NewRunner() *Runner {
	return &Runner{
		runtime: &runtime{},
	}
}

func (r *Runner) Init(mem uint, in io.ByteReader, stdout io.Writer, prog *ast.Program) {
	r.prog = prog
	r.runtime.init(mem, in, stdout)
}

func (r *Runner) Run() {
	ast.Walk(r.runtime, r.prog)
}
