package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/alienvspredator/brainfuck-yacc/pkg/brainfuck"
)

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

		brainfuck.Run(line)
	}
}
