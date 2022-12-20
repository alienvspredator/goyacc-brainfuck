package parser

import (
	"github.com/alienvspredator/brainfuck-yacc/pkg/brainfuck/ast"
	"github.com/alienvspredator/brainfuck-yacc/pkg/brainfuck/scanner"
	"github.com/alienvspredator/brainfuck-yacc/pkg/brainfuck/token"
)

type parser struct {
	source  *token.Source
	errors  scanner.ErrorList
	scanner scanner.Scanner

	pos token.Pos
	tok token.Token
	lit string

	// nestLev is used to track and limit the recursion depth
	// during parsing.
	nestLev int
}

func (p *parser) init(src []byte) {
	p.source = token.NewSource(len(src))

	eh := func(pos token.Position, msg string) { p.errors.Add(pos, msg) }
	p.scanner.Init(p.source, src, eh)

	p.next()
}

const maxNestLev int = 1e5

func incNestLev(p *parser) *parser {
	p.nestLev++
	if p.nestLev > maxNestLev {
		p.error(p.pos, "exceeded max nesting depth")
		panic(bailout{})
	}

	return p
}

func decNestLev(p *parser) {
	p.nestLev--
}

func (p *parser) next() {
	p.pos, p.tok, p.lit = p.scanner.Scan()
}

type bailout struct {
	pos token.Pos
	msg string
}

func (p *parser) error(pos token.Pos, msg string) {
	p.errors.Add(p.source.Position(pos), msg)
}

func (p *parser) errorExpected(pos token.Pos, msg string) {
	msg = "expected " + msg
	if pos == p.pos {
		msg += ", found '" + p.tok.String() + "'"
	}

	p.error(pos, msg)
}

func (p *parser) expect(tok token.Token) token.Pos {
	pos := p.pos
	if p.tok != tok {
		p.errorExpected(pos, "'"+tok.String()+"'")
	}
	p.next() // make progress
	return pos
}

func (p *parser) expect2(tok token.Token) (pos token.Pos) {
	if p.tok == tok {
		pos = p.pos
	} else {
		p.errorExpected(p.pos, "'"+tok.String()+"'")
	}
	p.next() // make progress
	return
}

func (p *parser) advance() {
	for ; p.tok != token.EOF; p.next() {
	}
}

func (p *parser) parseIncPtr() *ast.IncPtr {
	pos := p.expect(token.IncPtr)
	return &ast.IncPtr{
		Inc: pos,
	}
}

func (p *parser) parseDecPtr() *ast.DecPtr {
	pos := p.expect(token.DecPtr)
	return &ast.DecPtr{
		Dec: pos,
	}
}

func (p *parser) parseIncByte() *ast.IncByte {
	pos := p.expect(token.IncByte)
	return &ast.IncByte{
		Inc: pos,
	}
}

func (p *parser) parseDecByte() *ast.DecByte {
	pos := p.expect(token.DecByte)
	return &ast.DecByte{
		Dec: pos,
	}
}

func (p *parser) parseOutputByte() *ast.OutputByte {
	pos := p.expect(token.OutputByte)
	return &ast.OutputByte{
		TokenPos: pos,
	}
}

func (p *parser) parseInputByte() *ast.InputByte {
	pos := p.expect(token.InputByte)
	return &ast.InputByte{
		TokenPos: pos,
	}
}

func (p *parser) parseNodeList() (list []ast.Node) {
	for p.tok != token.LoopClose && p.tok != token.EOF {
		list = append(list, p.parseNode())
	}

	return
}

func (p *parser) parseBody() *ast.Body {
	return &ast.Body{
		Body: p.pos,
		List: p.parseNodeList(),
	}
}

func (p *parser) parseLoop() *ast.Loop {
	pos := p.pos
	p.expect(token.LoopOpen)
	body := p.parseBody()
	p.expect2(token.LoopClose)

	return &ast.Loop{
		Loop: pos,
		Body: body,
	}
}

func (p *parser) parseNode() ast.Node {
	defer decNestLev(incNestLev(p))

	switch p.tok {
	case token.IncPtr:
		return p.parseIncPtr()
	case token.DecPtr:
		return p.parseDecPtr()
	case token.IncByte:
		return p.parseIncByte()
	case token.DecByte:
		return p.parseDecByte()
	case token.OutputByte:
		return p.parseOutputByte()
	case token.InputByte:
		return p.parseInputByte()
	case token.LoopOpen:
		return p.parseLoop()
	default:
		pos := p.pos
		p.errorExpected(pos, "node")
		p.advance()
		return &ast.BadNode{From: pos, To: p.pos}
	}
}

func (p *parser) parseProgram() *ast.Program {
	if p.errors.Len() != 0 {
		return nil
	}

	var nodes []ast.Node
	for p.tok != token.EOF {
		nodes = append(nodes, p.parseNode())
	}

	return &ast.Program{
		Body: &ast.Body{
			List: nodes,
		},
	}
}
