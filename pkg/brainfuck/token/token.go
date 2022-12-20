package token

import (
	"strconv"
)

type Token int

const (
	ILLEGAL Token = iota
	EOF
	IncPtr     // <
	DecPtr     // >
	IncByte    // +
	DecByte    // -
	OutputByte // .
	InputByte  // ,
	LoopOpen   // [
	LoopClose  // ]
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF: "EOF",

	IncPtr: "<",
	DecPtr: ">",

	IncByte: "+",
	DecByte: "-",

	OutputByte: ".",
	InputByte:  ",",

	LoopOpen:  "[",
	LoopClose: "]",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}

	return s
}
