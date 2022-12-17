// Package ast declares brainfuck AST.
package ast

import (
	"github.com/alienvspredator/brainfuck-yacc/pkg/brainfuck/token"
)

type Node interface {
	Pos() token.Pos
	End() token.Pos
}

type Program struct {
	Body *Body
}

func (p *Program) Pos() token.Pos {
	return 0
}

func (p *Program) End() token.Pos {
	return 0
}

type Visitor interface {
	Visit(node Node) Visitor
}

type IncPtr struct {
	Inc token.Pos
}

func (x *IncPtr) Pos() token.Pos {
	return x.Inc
}

func (x *IncPtr) End() token.Pos {
	return x.Inc + 1
}

type DecPtr struct {
	Dec token.Pos
}

func (x *DecPtr) Pos() token.Pos {
	return x.Dec
}

func (x *DecPtr) End() token.Pos {
	return x.Dec + 1
}

type IncByte struct {
	Inc token.Pos
}

func (x *IncByte) Pos() token.Pos {
	return x.Inc
}

func (x *IncByte) End() token.Pos {
	return x.Inc + 1
}

type DecByte struct {
	Dec token.Pos
}

func (x *DecByte) Pos() token.Pos {
	return x.Dec
}

func (x *DecByte) End() token.Pos {
	return x.Dec + 1
}

type OutputByte struct {
	TokenPos token.Pos
}

func (x *OutputByte) Pos() token.Pos {
	return x.TokenPos
}

func (x *OutputByte) End() token.Pos {
	return x.TokenPos + 1
}

type InputByte struct {
	TokenPos token.Pos
}

func (x *InputByte) Pos() token.Pos {
	return x.TokenPos
}

func (x *InputByte) End() token.Pos {
	return x.TokenPos + 1
}

type Loop struct {
	Loop token.Pos
	Body *Body
}

func (x *Loop) Pos() token.Pos {
	return x.Loop
}

func (x *Loop) End() token.Pos {
	return x.Body.End()
}

type Body struct {
	Body token.Pos
	List []Node
}

func (x *Body) Pos() token.Pos {
	return x.Body
}

func (x *Body) End() token.Pos {
	if n := len(x.List); n > 0 {
		return x.List[n-1].End()
	}

	return x.Body + 1
}
