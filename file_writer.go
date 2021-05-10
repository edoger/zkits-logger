// Copyright 2021 The ZKits Project Authors.
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
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/edoger/zkits-logger/internal"
)

// NewFileWriter creates and returns an io.WriteCloser instance from the given path.
// The max parameter is used to limit the maximum size of the log file, if it is 0, the
// file size limit will be disabled. The log files that exceed the maximum size limit
// will be renamed.
func NewFileWriter(name string, max uint32) (io.WriteCloser, error) {
	if abs, err := filepath.Abs(name); err != nil {
		return nil, err
	} else {
		name = abs
	}
	// Before opening the log file, make sure that the directory must exist.
	if err := os.MkdirAll(filepath.Dir(name), 0766); err != nil {
		return nil, err
	}
	fd, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	i, err := fd.Stat()
	if err != nil {
		_ = fd.Close()
		return nil, err
	}
	return &fileWriter{name: name, max: max, size: uint32(i.Size()), fd: fd}, nil
}

// MustNewFileWriter is like NewFileWriter, but triggers a panic when an error occurs.
func MustNewFileWriter(name string, max uint32) io.WriteCloser {
	w, err := NewFileWriter(name, max)
	if err != nil {
		panic(err)
	}
	return w
}

// The built-in log file writer.
type fileWriter struct {
	name string
	mu   sync.RWMutex
	max  uint32
	size uint32
	fd   *os.File
}

// Write is an implementation of the io.WriteCloser interface, used to write a single
// log data to a file.
func (w *fileWriter) Write(b []byte) (n int, err error) {
	w.mu.RLock()
	n, err = w.fd.Write(b)
	w.mu.RUnlock()
	if n > 0 && w.max != 0 && atomic.AddUint32(&w.size, uint32(n)) >= w.max {
		w.swap()
	}
	return
}

// Switch to the new log file.
func (w *fileWriter) swap() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if atomic.LoadUint32(&w.size) < w.max {
		return
	}
	// We need to make sure that the log file is written to disk before closing the file.
	if err := w.fd.Sync(); err != nil {
		internal.EchoError("Failed to sync log file %s: %s.", w.name, err)
		return
	}
	if err := w.fd.Close(); err != nil {
		// We need to ignore the error that the file is closed.
		if e, ok := err.(*os.PathError); !ok || e.Err != os.ErrClosed {
			internal.EchoError("Failed to close log file %s: %s.", w.name, err)
			return
		}
	}
	// We use a second-level date as the suffix name of the archive log,
	// which may change in the future.
	suffix := time.Now().Format("20060102150405")
	if err := os.Rename(w.name, w.name+"."+suffix); err != nil {
		internal.EchoError("Failed to rename log file %s to %s.%s: %s.", w.name, w.name, suffix, err)
		return
	}
	fd, err := os.OpenFile(w.name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		internal.EchoError("Failed to open log file %s: %s.", w.name, err)
		return
	}
	w.fd = fd
	atomic.StoreUint32(&w.size, 0)
}

// Close is an implementation of the io.WriteCloser interface.
func (w *fileWriter) Close() (err error) {
	w.mu.RLock()
	err = w.fd.Close()
	w.mu.RUnlock()
	return
}
