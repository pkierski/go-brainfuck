// Package brainfuck implementation, just for training purposes.
//
// This implementation allows to execute Brainfuck programs
// step by step and examining (and altering) Brainfuck machine state
// (memory, instruction pointer and memory pointer).
package brainfuck

import (
	"errors"
	"io"
)

// Brainfuck machine state.
type Brainfuck struct {
	// Code to execute.
	Code string
	// Instruction pointer.
	IP int
	// Memory.
	Memory []byte
	// Memory pointer.
	MemPtr   int
	input    io.Reader
	output   io.Writer
	ioBuffer []byte
}

// New creates Brainfuck machine state.
//
// Instruction pointer is set on first instruction in code
//
// Memory is filled with zeros and memory pointer is set to 0 (first cell in array).
//
// input and output may be empty if code doesn't use any output or input operation.
func New(code string, memSize int, input io.Reader, output io.Writer) *Brainfuck {
	return &Brainfuck{
		Code:     code,
		Memory:   make([]byte, memSize),
		input:    input,
		output:   output,
		ioBuffer: make([]byte, 1),
	}
}

// Step executes one step on current state.
//
// Returns Finished() and error (if encountered).
//
// Possible sources of errors: I/O errors (errors on '.' or ',' command or
// nil as input or output when state was created), unbalanced loops ('[' without
// ']' or vice-versa), corrupted state (negative instruction pointer).
//
// In case of error instruction pointer remains at instruction causing error.
func (bf *Brainfuck) Step() (bool, error) {
	if bf.IP < 0 {
		return true, errors.New("invalid instruction pointer")
	}

	if bf.Finished() {
		return true, errors.New("already finished")
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
			return bf.Finished(), errors.New("no output defined on output operation")
		}
		_, err := bf.output.Write(bf.Memory[bf.MemPtr : bf.MemPtr+1])
		if err != nil {
			bf.IP++
			return bf.Finished(), err
		}
	case ',':
		if bf.input == nil {
			return bf.Finished(), errors.New("no input defined on input operation")
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
					return true, errors.New("loop not closed")
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
					return true, errors.New("no matching opening loop instruction")
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

// Finished returns true if instruction pointer reached end of the code.
func (bf *Brainfuck) Finished() bool {
	return bf.IP >= len(bf.Code)
}

// Run executes program until Finished() is false or error encountered.
//
// In case of error execution is stopped and error is returned.
// Instruction pointer remains on instruction causing error.
func (bf *Brainfuck) Run() error {
	var err error
	var finished bool
	for finished, err = bf.Step(); !finished && err == nil; finished, err = bf.Step() {
		// intentionally empty, all work are done in Step()
	}
	return err
}
