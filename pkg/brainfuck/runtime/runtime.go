// Package runtime is a runtime for brainfuck.
package runtime

import (
	"fmt"
	"io"

	"github.com/alienvspredator/brainfuck-yacc/pkg/brainfuck/ast"
)

type runtime struct {
	stdin  io.ByteReader
	stdout io.Writer

	ptr    uint
	memory []rune
}

func (r *runtime) init(mem uint, stdin io.ByteReader, stdout io.Writer) {
	r.stdin = stdin
	r.stdout = stdout

	r.ptr = 0
	r.memory = make([]rune, mem)
}

func (r *runtime) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Loop:
		return r.visitLoop(n)
	case *ast.IncPtr:
		return r.visitIncPtr(n)
	case *ast.DecPtr:
		return r.visitDecPtr(n)
	case *ast.IncByte:
		return r.visitIncByte(n)
	case *ast.DecByte:
		return r.visitDecByte(n)
	case *ast.OutputByte:
		return r.visitOutputByte(n)
	case *ast.InputByte:
		return r.visitInputByte(n)
	}

	return r
}

func (r *runtime) currPtrVal() rune {
	return r.memory[r.ptr]
}

func (r *runtime) visitLoop(n *ast.Loop) ast.Visitor {
	for r.currPtrVal() != 0 {
		ast.Walk(r, n.Body)
	}

	return nil
}

func (r *runtime) visitIncPtr(_ *ast.IncPtr) ast.Visitor {
	r.ptr++
	return r
}

func (r *runtime) visitDecPtr(_ *ast.DecPtr) ast.Visitor {
	r.ptr--
	return r
}

func (r *runtime) visitIncByte(_ *ast.IncByte) ast.Visitor {
	r.memory[r.ptr]++
	return r
}

func (r *runtime) visitDecByte(_ *ast.DecByte) ast.Visitor {
	r.memory[r.ptr]--
	return r
}

func (r *runtime) visitOutputByte(_ *ast.OutputByte) ast.Visitor {
	fmt.Fprintf(r.stdout, "%c", r.currPtrVal())
	return r
}

func (r *runtime) visitInputByte(_ *ast.InputByte) ast.Visitor {
	b, err := r.stdin.ReadByte()
	if err != nil {
		panic(err)
	}

	r.memory[r.ptr] = rune(b)
	return r
}
