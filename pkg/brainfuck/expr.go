package brainfuck

import (
	"fmt"
)

type Expr interface {
	Eval(state *State)
}

type addExpr struct{}

func NewAdd() Expr {
	return addExpr{}
}

func (_ addExpr) Eval(state *State) {
	state.Memory[state.Ptr]++
}

type subExpr struct{}

func NewSubtract() Expr {
	return subExpr{}
}

func (_ subExpr) Eval(state *State) {
	state.Memory[state.Ptr]--
}

type incPtrExpr struct{}

func NewIncPtr() Expr {
	return incPtrExpr{}
}
func (_ incPtrExpr) Eval(state *State) {
	state.Ptr++
}

type decPtrExpr struct{}

func NewDecPtr() Expr {
	return decPtrExpr{}
}

func (_ decPtrExpr) Eval(state *State) {
	state.Ptr--
}

type loopExpr struct {
	expr Expr
}

func NewLoop(expr Expr) Expr {
	return loopExpr{expr: expr}
}

func (l loopExpr) Eval(state *State) {
	for state.Memory[state.Ptr] != 0 {
		l.expr.Eval(state)
	}
}

type exprList []Expr

func NewList(expr ...Expr) Expr {
	return exprList(expr)
}

func (l exprList) Eval(state *State) {
	for _, expr := range l {
		expr.Eval(state)
	}
}

type outputExpr struct{}

func NewOutput() Expr {
	return outputExpr{}
}

func (_ outputExpr) Eval(state *State) {
	fmt.Printf("%d ", state.Memory[state.Ptr])
}

type inputExpr struct{}

func NewInput() Expr {
	return inputExpr{}
}

func (_ inputExpr) Eval(state *State) {
	fmt.Println("input expression isn't implemented")
}
