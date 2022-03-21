// Copyright 2022 The ZKits Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"bytes"
	"sync"
	"testing"
)

func TestNewMutexWriter(t *testing.T) {
	w := new(bytes.Buffer)
	if got := NewMutexWriter(w); got == nil {
		t.Fatal("NewMutexWriter(): return nil")
	}
}

type testUnsafeWriter struct {
	n uint32
	m map[uint32][]byte
}

func (w *testUnsafeWriter) Write(p []byte) (int, error) {
	w.n++
	p2 := make([]byte, len(p))
	copy(p2, p)
	w.m[w.n] = p2
	return len(p2), nil
}

func TestMutexWriter_Write(t *testing.T) {
	iw := &testUnsafeWriter{
		m: make(map[uint32][]byte),
	}
	w := NewMutexWriter(iw)
	l := New("test").SetOutput(w)
	wg := new(sync.WaitGroup)
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			l.Info("test")
		}()
	}
	wg.Wait()
	if len(iw.m) != 10 {
		t.Fatalf("MutexWriter: want %d, got %d", 10, len(iw.m))
	}
	if iw.n != 10 {
		t.Fatalf("MutexWriter: want %d, got %d", 10, iw.n)
	}
}
