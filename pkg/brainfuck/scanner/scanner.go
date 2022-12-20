package scanner

import (
	"fmt"
	"unicode/utf8"

	"github.com/alienvspredator/brainfuck-yacc/pkg/brainfuck/token"
)

type ErrorHandler func(pos token.Position, msg string)

type Scanner struct {
	// immutable state
	source *token.Source
	src    []byte
	err    ErrorHandler

	// scanning state
	ch         rune // current character
	offset     int  // character offset
	rdOffset   int  // reading offset (position after current character)
	lineOffset int  // current line offset

	// public state
	ErrorCount int
}

const (
	bom = 0xFEFF
	eof = -1
)

// next reads the next Unicode char into s.ch. s.ch < 0 means eof.
func (s *Scanner) next() {
	if s.rdOffset < len(s.src) {
		s.offset = s.rdOffset
		if s.ch == '\n' {
			s.lineOffset = s.offset
		}
		r, w := rune(s.src[s.rdOffset]), 1
		switch {
		case r == 0:
			s.error(s.offset, "illegal character NUL")
		case r >= utf8.RuneSelf:
			r, w = utf8.DecodeRune(s.src[s.rdOffset:])
			if r == utf8.RuneError && w == 1 {
				s.error(s.offset, "illegal UTF-8 encoding")
			} else if r == bom && s.offset > 0 {
				s.error(s.offset, "illegal byte order mark")
			}
		}

		s.rdOffset += w
		s.ch = r
	} else {
		s.offset = len(s.src)
		if s.ch == '\n' {
			s.lineOffset = s.offset
		}
		s.ch = eof
	}
}

func (s *Scanner) peek() byte {
	if s.rdOffset < len(s.src) {
		return s.src[s.rdOffset]
	}

	return 0
}

func (s *Scanner) Init(source *token.Source, src []byte, err ErrorHandler) {
	if source.Size() != len(src) {
		panic(fmt.Sprintf("source size (%d) does not match src len (%d)", source.Size(), len(src)))
	}

	s.source = source
	s.src = src
	s.err = err

	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.lineOffset = 0
	s.ErrorCount = 0

	s.next()
	if s.ch == bom {
		s.next()
	}
}

func (s *Scanner) error(offs int, msg string) {
	if s.err != nil {
		s.err(s.source.Position(token.Pos(offs)), msg)
	}
	s.ErrorCount++
}

func (s *Scanner) errorf(offs int, format string, args ...any) {
	s.error(offs, fmt.Sprintf(format, args...))
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

func (s *Scanner) Scan() (pos token.Pos, tok token.Token, lit string) {
	s.skipWhitespace()

	pos = token.Pos(s.offset)
	ch := s.ch

	s.next()
	switch ch {
	case -1:
		tok = token.EOF
	case '+':
		tok = token.IncByte
	case '-':
		tok = token.DecByte
	case '>':
		tok = token.IncPtr
	case '<':
		tok = token.DecPtr
	case '[':
		tok = token.LoopOpen
	case ']':
		tok = token.LoopClose
	case '.':
		tok = token.OutputByte
	case ',':
		tok = token.InputByte
	default:
		if ch != bom {
			s.errorf(int(pos), "illegal character %#U", ch)
		}
		tok = token.ILLEGAL
		lit = string(ch)
	}

	return
}
