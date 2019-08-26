package brainfuck

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyProgram(t *testing.T) {
	bf := New("", 100, nil, nil)
	err := bf.Run()

	assert.Nil(t, err)
	assert.Equal(t, bf.IP, 0)
}

func TestHelloWorld(t *testing.T) {
	output := &strings.Builder{}
	bf := New(
		"+[-[<<[+[--->]-[<<<]]]>>>-]>-.---.>..>.<<<<-.<+.>>>>>.>.<<.<-.",
		100,
		nil,
		output)
	err := bf.Run()
	assert.Nil(t, err)
	assert.Equal(t, output.String(), "hello world")
}
