package token

import (
	"fmt"
	"sync"
)

type Position struct {
	Offset int
	Line   int
	Column int
}

func (pos *Position) IsValid() bool {
	return pos.Line > 0
}

func (pos *Position) String() (s string) {
	if pos.IsValid() {
		s += fmt.Sprintf("%d", pos.Line)
		if pos.Column != 0 {
			s += fmt.Sprintf(":%d", pos.Column)
		}
	}

	if s == "" {
		s = "-"
	}

	return
}

type Pos int

const NoPos Pos = 0

func (p Pos) IsValid() bool {
	return p != NoPos
}

type Source struct {
	size int

	// lines and infos are protected by mutex
	mutex sync.Mutex
	lines []int // lines contains the offset of the first character for each line (the first entry is always 0)
}

func NewSource(size int) *Source {
	return &Source{size: size, lines: []int{0}}
}

func (s *Source) Size() int {
	return s.size
}

func (s *Source) LineCount() int {
	s.mutex.Lock()
	n := len(s.lines)
	s.mutex.Unlock()
	return n
}

func (s *Source) AddLine(offset int) {
	s.mutex.Lock()
	if i := len(s.lines); (i == 0 || s.lines[i-1] < offset) && offset < s.size {
		s.lines = append(s.lines, offset)
	}

	s.mutex.Unlock()
}

func (s *Source) MergeLine(line int) {
	if line < 1 {
		panic(fmt.Sprintf("invalid line number %d (should be >= 1)", line))
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	if line >= len(s.lines) {
		panic(fmt.Sprintf("invalid line number %d (should be < %d)", line, len(s.lines)))
	}

	copy(s.lines[line:], s.lines[line+1:])
	s.lines = s.lines[:len(s.lines)-1]
}

func (s *Source) SetLines(lines []int) bool {
	// verify validity of lines table
	size := s.size
	for i, offset := range lines {
		if i > 0 && offset <= lines[i-1] || size <= offset {
			return false
		}
	}

	// set lines table
	s.mutex.Lock()
	s.lines = lines
	s.mutex.Unlock()
	return true
}

func (s *Source) SetLinesForContent(content []byte) {
	var lines []int
	line := 0
	for offset, b := range content {
		if line >= 0 {
			lines = append(lines, line)
		}
		line = -1
		if b == '\n' {
			line = offset + 1
		}
	}

	// set lines table
	s.mutex.Lock()
	s.lines = lines
	s.mutex.Unlock()
}

func (s *Source) LineStart(line int) Pos {
	if line < 1 {
		panic(fmt.Sprintf("invalid line number %d (should be >= 1)", line))
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if line > len(s.lines) {
		panic(fmt.Sprintf("invalid line number %d (should be < %d)", line, len(s.lines)))
	}

	return Pos(s.lines[line-1])
}

// Line returns the line number for the given file position p;
// p must be a Pos value in that file or NoPos.
func (s *Source) Line(p Pos) int {
	return s.Position(p).Line
}

func (s *Source) unpack(offset int) (line, column int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if i := searchInts(s.lines, offset); i >= 0 {
		line, column = i+1, offset-s.lines[i]
	}

	return
}

func (s *Source) position(p Pos) (pos Position) {
	offset := int(p)
	pos.Offset = offset
	pos.Line, pos.Column = s.unpack(offset)
	return
}

func (s *Source) Position(p Pos) (pos Position) {
	if p != NoPos {
		pos = s.position(p)
	}

	return
}

func searchInts(a []int, x int) int {
	i, j := 0, len(a)
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if a[h] <= x {
			i = h + 1
		} else {
			j = h
		}
	}
	return i - 1
}
