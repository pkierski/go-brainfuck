package brainfuck_test

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/pkierski/go-brainfuck"
	"github.com/stretchr/testify/assert"
)

func TestEmptyProgram(t *testing.T) {
	bf := brainfuck.New("", 100, nil, nil)
	err := bf.Run()

	assert.NotNil(t, err)
	assert.Equal(t, bf.IP, 0)
}

func TestCorruptedIP(t *testing.T) {
	bf := brainfuck.New("", 100, nil, nil)
	bf.IP = -2
	err := bf.Run()

	assert.NotNil(t, err)
}

func Example(t *testing.T) {
	output := &strings.Builder{}
	bf := brainfuck.New(
		"+[-[<<[+[--->]-[<<<]]]>>>-]>-.---.>..>.<<<<-.<+.>>>>>.>.<<.<-.",
		100,
		nil,
		output)
	err := bf.Run()
	if err == nil {
		fmt.Println(output.String())
	} else {
		panic(err)
	}
	// Output:
	// hello world
}

func TestMemPtrWrapping(t *testing.T) {
	bf := brainfuck.New("<>", 100, nil, nil)
	var end bool
	var err error

	end, err = bf.Step()
	assert.Nil(t, err)
	assert.False(t, end)
	assert.Equal(t, bf.MemPtr, len(bf.Memory)-1)

	end, err = bf.Step()
	assert.Nil(t, err)
	assert.True(t, end)
	assert.Equal(t, bf.MemPtr, 0)
}

func TestUnbalancedOpenLoop(t *testing.T) {
	bf := brainfuck.New("[[-]", 100, nil, nil)
	err := bf.Run()
	assert.NotNil(t, err)
}

func TestUnbalancedCloseLoop(t *testing.T) {
	bf := brainfuck.New("[-]+]", 100, nil, nil)
	err := bf.Run()
	assert.NotNil(t, err)
}

func TestNoOutputOnDot(t *testing.T) {
	bf := brainfuck.New(".", 100, nil, nil)
	err := bf.Run()
	assert.NotNil(t, err)
}

type errorOutput struct{}

func (*errorOutput) Write([]byte) (int, error) {
	return 0, errors.New("Cannot write")
}

func TestNoErrorOutputOnDot(t *testing.T) {
	bf := brainfuck.New(".", 100, nil, &errorOutput{})
	err := bf.Run()
	assert.NotNil(t, err)
}

func TestNoInputOnComma(t *testing.T) {
	bf := brainfuck.New(",", 100, nil, nil)
	err := bf.Run()
	assert.NotNil(t, err)
}

func TestEmptyInputOnComma(t *testing.T) {
	emptyInput := bytes.NewBuffer([]byte{})
	bf := brainfuck.New("+,", 100, emptyInput, nil)
	err := bf.Run()
	assert.Nil(t, err)
	assert.Equal(t, bf.Memory[0], byte(0))
}

type errorInput struct{}

func (*errorInput) Read([]byte) (int, error) {
	return 0, errors.New("Cannot read")
}

func TestErrorInputOnComma(t *testing.T) {
	errorInput := &errorInput{}
	bf := brainfuck.New("+,", 100, errorInput, nil)
	err := bf.Run()
	assert.NotNil(t, err)
	assert.Equal(t, bf.Memory[0], byte(1))
}
