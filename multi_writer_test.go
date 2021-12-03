package logger

import (
	"bytes"
	"errors"
	"testing"
)

func TestNewMultiWriter(t *testing.T) {
	if w := NewMultiWriter(); w == nil {
		t.Fatal("NewMultiWriter(): return nil")
	}
	if w := NewMultiWriter(bytes.NewBuffer(nil)); w == nil {
		t.Fatal("NewMultiWriter(): return nil")
	}
	if w := NewMultiWriter(NewMultiWriter(), NewMultiWriter(bytes.NewBuffer(nil))); w == nil {
		t.Fatal("NewMultiWriter(): return nil")
	}
}

func TestMultiWriter_Write(t *testing.T) {
	b1 := bytes.NewBufferString("b1")
	b2 := bytes.NewBufferString("b2")
	mw := NewMultiWriter(b1, b2)
	if n, err := mw.Write([]byte("test")); err != nil {
		t.Fatalf("NewMultiWriter().Write(): %s", err)
	} else {
		if n != 4 {
			t.Fatalf("NewMultiWriter().Write(): return %d", n)
		}
	}
	if got := b1.String(); got != "b1test" {
		t.Fatalf("NewMultiWriter().Write(): got %s", got)
	}
	if got := b2.String(); got != "b2test" {
		t.Fatalf("NewMultiWriter().Write(): got %s", got)
	}
}

func TestMultiWriter_Write_WithoutWriter(t *testing.T) {
	w := NewMultiWriter()
	if n, err := w.Write([]byte("test")); err != nil {
		t.Fatalf("NewMultiWriter().Write(): %s", err)
	} else {
		if n != 4 {
			t.Fatalf("NewMultiWriter().Write(): return %d", n)
		}
	}
}

type testFixedReturnValueWriter struct {
	n   int
	err error
}

func (w *testFixedReturnValueWriter) Write(p []byte) (int, error) {
	return w.n, w.err
}

func TestMultiWriter_Write_WithErrorWriter(t *testing.T) {
	w1 := new(testFixedReturnValueWriter)
	w2 := bytes.NewBufferString("w2")
	w := NewMultiWriter(w1, w2)

	if n, err := w.Write([]byte("test")); err == nil {
		t.Fatal("NewMultiWriter().Write(): return nil error")
	} else {
		if n != 4 {
			t.Fatalf("NewMultiWriter().Write(): return %d", n)
		}
	}
	if got := w2.String(); got != "w2test" {
		t.Fatalf("NewMultiWriter().Write(): got %s", got)
	}

	w2.Reset()
	w1.n = 4
	w1.err = errors.New("error")
	if n, err := w.Write([]byte("test")); err == nil {
		t.Fatal("NewMultiWriter().Write(): return nil error")
	} else {
		if n != 4 {
			t.Fatalf("NewMultiWriter().Write(): return %d", n)
		}
	}
	if got := w2.String(); got != "test" {
		t.Fatalf("NewMultiWriter().Write(): got %s", got)
	}
}
