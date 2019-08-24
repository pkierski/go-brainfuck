package brainfuck

import (
	"strings"
	"testing"
)

func TestEmptyProgram(t *testing.T) {
	bf := New("", 100, nil, nil)
	err := bf.Run()
	if err != nil {
		t.Errorf("Unexpected error=%v", err)
	}
	if bf.IP != 0 {
		t.Errorf("Unexpected IP value=%v", bf.IP)
	}
}

func TestHelloWorld(t *testing.T) {
	output := &strings.Builder{}
	bf := New(
		"+[-[<<[+[--->]-[<<<]]]>>>-]>-.---.>..>.<<<<-.<+.>>>>>.>.<<.<-.",
		100,
		nil,
		output)
	err := bf.Run()
	if err != nil {
		t.Errorf("Unexpected error=%v", err)
	}
	if output.String() != "hello world" {
		t.Errorf("Unexpected output: '%v'", output.String())
	}
}
