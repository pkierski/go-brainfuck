package brainfuck

import (
	"io"
)

type Brainfuck struct {
	Code     string
	IP       int
	Memory   []byte
	MemPtr   int
	input    io.Reader
	output   io.Writer
	ioBuffer []byte
}

func New(code string, memSize int, input io.Reader, output io.Writer) *Brainfuck {
	return &Brainfuck{
		Code:     code,
		Memory:   make([]byte, memSize),
		input:    input,
		output:   output,
		ioBuffer: make([]byte, 1),
	}
}

func (bf *Brainfuck) Step() (bool, error) {
	return bf.Finished(), nil
}

func (bf *Brainfuck) Finished() bool {
	return bf.IP >= len(bf.Code)
}

func (bf *Brainfuck) Run() error {
	var err error
	for finished, err := bf.Step(); !finished && err == nil; finished, err = bf.Step() {
	}
	return err
}
