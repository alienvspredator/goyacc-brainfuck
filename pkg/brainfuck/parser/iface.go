package parser

import (
	"bytes"
	"errors"
	"io"

	"github.com/alienvspredator/brainfuck-yacc/pkg/brainfuck/ast"
)

func ParseProgram(src any) (prog *ast.Program, err error) {
	text, err := readSource(src)
	if err != nil {
		return nil, err
	}

	var p parser
	defer func() {
		if e := recover(); e != nil {
			bail, ok := e.(bailout)
			if !ok {
				panic(e)
			} else if bail.msg != "" {
				p.errors.Add(p.source.Position(bail.pos), bail.msg)
			}
		}

		if prog == nil {
			prog = &ast.Program{
				Body: &ast.Body{},
			}
		}

		p.errors.Sort()
		err = p.errors.Err()
	}()

	p.init(text)
	prog = p.parseProgram()

	return
}

func readSource(src any) ([]byte, error) {
	switch s := src.(type) {
	case string:
		return []byte(s), nil
	case []byte:
		return s, nil
	case *bytes.Buffer:
		if s != nil {
			return s.Bytes(), nil
		}
	case io.Reader:
		return io.ReadAll(s)
	}

	return nil, errors.New("invalid source")
}
