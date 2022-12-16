package brainfuck

type Program struct {
	State *State
	Expr  Expr
}

func NewProgram(mem uint, expr Expr) *Program {
	return &Program{
		State: NewState(mem),
		Expr:  expr,
	}
}

func (p *Program) Run() {
	if p.Expr != nil {
		p.Expr.Eval(p.State)
	}
}

type State struct {
	Memory []uint16
	Ptr    uint
}

func NewState(mem uint) *State {
	return &State{Memory: make([]uint16, mem)}
}
