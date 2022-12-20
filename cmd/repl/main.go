package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"

	"github.com/alienvspredator/brainfuck-yacc/pkg/brainfuck/parser"
	"github.com/alienvspredator/brainfuck-yacc/pkg/brainfuck/runtime"
)

type fakeReader struct{}

func (f fakeReader) ReadByte() (byte, error) {
	return byte('A'), nil
}

func main() {
	in := bufio.NewReader(os.Stdin)

	for {
		if _, err := os.Stdout.WriteString("> "); err != nil {
			log.Fatalf("WriteString: %s", err)
		}

		line, err := in.ReadBytes('\n')
		if err == io.EOF {
			return
		}

		if err != nil {
			log.Fatalf("ReadBytes: %s", err)
		}

		prog, err := parser.ParseProgram(line)
		if err != nil {
			log.Printf("faield to parse the input: %s\n", err)
			continue
		}

		runner := runtime.NewRunner()

		var buf bytes.Buffer

		runner.Init(3000, fakeReader{}, &buf, prog)
		runner.Run()

		if buf.Len() > 0 {
			buf.WriteRune('\n')
			buf.WriteTo(os.Stdout)
		}
	}
}
