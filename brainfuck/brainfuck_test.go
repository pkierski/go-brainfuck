package brainfuck

import (
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
