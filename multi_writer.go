package logger

import (
	"io"
)

// NewMultiWriter creates a writer that duplicates its writes to all the provided writers.
// If any writer fails to write, it will continue to try to write to other writers.
// If you need to interrupt writing when any writer fails to write, please use io.MultiWriter.
func NewMultiWriter(writers ...io.Writer) io.Writer {
	ws := make([]io.Writer, 0, len(writers))
	for i, j := 0, len(writers); i < j; i++ {
		if w, ok := writers[i].(*multiWriter); ok {
			ws = append(ws, w.writers...)
		} else {
			ws = append(ws, writers[i])
		}
	}
	return &multiWriter{ws}
}

// This is a multi-channel writer.
type multiWriter struct {
	writers []io.Writer
}

// Write is the implementation of io.Writer interface.
// If there are no writers available, we always return success. We only return the first
// error encountered. We only return the maximum number of bytes written.
func (w *multiWriter) Write(p []byte) (int, error) {
	if j := len(w.writers); j > 0 {
		var err error
		var n int
		for i := 0; i < j; i++ {
			k, e := w.writers[i].Write(p)
			if e == nil && k != len(p) {
				e = io.ErrShortWrite
			}
			if err == nil {
				err = e
			}
			if k > n {
				n = k
			}
		}
		return n, err
	}
	return len(p), nil
}
