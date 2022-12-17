package token

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
