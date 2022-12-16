//go:generate goyacc -o lex.go -p "brainfuck" lex.y
package brainfuck

func Run(line []byte) int {
	return brainfuckParse(&brainfuckLex{line: line})
}
