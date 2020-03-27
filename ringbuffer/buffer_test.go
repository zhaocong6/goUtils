package ringbuffer

import (
	"testing"
)

func TestNew(t *testing.T) {
	b := New(Options{Capacity: 1})

	for i := 0; i < 1000000; i++ {
		b.Write(i)
	}
}
