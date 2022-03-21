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
	"io"
	"sync"
)

// NewMutexWriter creates and returns a mutex writer.
// Mutually exclusive writers ensure mutual exclusion of individual write calls.
func NewMutexWriter(w io.Writer) io.Writer {
	return &mutexWriter{w: w}
}

// This is an implementation of the built-in mutex writer.
type mutexWriter struct {
	w  io.Writer
	mu sync.Mutex
}

// Write is an implementation of the io.Writer interface.
func (w *mutexWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	n, err = w.w.Write(p)
	return
}
