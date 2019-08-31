package brainfuck

import (
	"errors"
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
	if bf.IP < 0 {
		return true, errors.New("Invalid instruction pointer")
	}

	if bf.Finished() {
		return true, errors.New("Already finished")
	}

	switch bf.Code[bf.IP] {
	case '+':
		bf.Memory[bf.MemPtr]++
	case '-':
		bf.Memory[bf.MemPtr]--
	case '<':
		bf.MemPtr = (bf.MemPtr + len(bf.Memory) - 1) % len(bf.Memory)
	case '>':
		bf.MemPtr = (bf.MemPtr + 1) % len(bf.Memory)
	case '.':
		if bf.output == nil {
			return bf.Finished(), errors.New("No output defined on output operation")
		}
		_, err := bf.output.Write(bf.Memory[bf.MemPtr : bf.MemPtr+1])
		if err != nil {
			bf.IP++
			return bf.Finished(), err
		}
	case ',':
		if bf.input == nil {
			return bf.Finished(), errors.New("No input defined on input operation")
		}
		_, err := bf.input.Read(bf.ioBuffer)
		if err != nil {
			bf.IP++
			if err == io.EOF {
				bf.Memory[bf.MemPtr] = 0
				break
			}
			return bf.Finished(), err
		}
	case '[':
		if bf.Memory[bf.MemPtr] == 0 {
			neestedLoops := 1
			newIP := bf.IP
			for {
				newIP++
				if newIP == len(bf.Code) {
					return true, errors.New("Loop not closed")
				}
				switch bf.Code[newIP] {
				case '[':
					neestedLoops++
				case ']':
					neestedLoops--
				}
				if neestedLoops == 0 {
					break
				}
			}
			bf.IP = newIP
		}
	case ']':
		if bf.Memory[bf.MemPtr] != 0 {
			newIP := bf.IP
			neestedLoops := 1
			for {
				newIP--
				if newIP == -1 {
					return true, errors.New("Unbalanced closing loop")
				}
				switch bf.Code[newIP] {
				case '[':
					neestedLoops--
				case ']':
					neestedLoops++
				}
				if neestedLoops == 0 {
					break
				}
			}
			bf.IP = newIP
		}
	}
	bf.IP++
	return bf.Finished(), nil
}

func (bf *Brainfuck) Finished() bool {
	return bf.IP >= len(bf.Code)
}

func (bf *Brainfuck) Run() error {
	var err error
	var finished bool
	for finished, err = bf.Step(); !finished && err == nil; finished, err = bf.Step() {
		// intentionally empty, all work are done in Step()
	}
	return err
}
