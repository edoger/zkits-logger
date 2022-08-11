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
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/edoger/zkits-logger/internal"
)

// NewFileWriter creates and returns an io.WriteCloser instance from the given path.
// The max parameter is used to limit the maximum size of the log file, if it is 0, the
// file size limit will be disabled.
// The log files that exceed the maximum size limit will be renamed.
// The backup parameter can limit the maximum number of backup log files retained.
// The log file writer we returned does not restrict concurrent writing. If necessary,
// you can use the writer wrapper with lock provided by us.
func NewFileWriter(name string, max, backup uint32) (io.WriteCloser, error) {
	if abs, err := filepath.Abs(name); err != nil {
		return nil, err
	} else {
		name = abs
	}
	// Before opening the log file, make sure that the directory must exist.
	// This requires the user who runs the program to have the authority to create the log
	// file directory.
	if err := os.MkdirAll(filepath.Dir(name), 0766); err != nil {
		return nil, err
	}
	w := &fileWriter{name: name, max: max, backup: backup}
	if fd, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		return nil, err
	} else {
		w.fd = fd
	}
	if max > 0 {
		// If the log file size is limited, then we need to determine the size of the
		// content in the current log file so that we can rotate the log file correctly.
		if i, err := w.fd.Stat(); err != nil {
			_ = w.fd.Close()
			return nil, err
		} else {
			w.size = uint32(i.Size())
		}
	}
	return w, nil
}

// MustNewFileWriter is like NewFileWriter, but triggers a panic when an error occurs.
func MustNewFileWriter(name string, max, backup uint32) io.WriteCloser {
	w, err := NewFileWriter(name, max, backup)
	if err != nil {
		panic(err)
	}
	return w
}

// The built-in log file writer.
type fileWriter struct {
	name   string
	max    uint32 // The maximum size of the log file.
	backup uint32 // The maximum number of backup log files.
	size   uint32 // The current log file size.
	mu     sync.RWMutex
	fd     *os.File
}

// Write is an implementation of the io.WriteCloser interface, used to write a single
// log data to a file.
func (w *fileWriter) Write(b []byte) (n int, err error) {
	w.mu.RLock()
	n, err = w.fd.Write(b)
	w.mu.RUnlock()
	// If the size of the log file is limited, we use write first and then rotate mode to
	// ensure write efficiency as much as possible, but this leads to the fact that the log
	// file writes some extra parts in the rotation gap when the log volume is very large.
	// For log data, the size of the log file that is finally retained will be larger than
	// the limited size. However, I think that higher log writing efficiency is more valuable
	// than more precise log file size. Therefore, if you really need a very strict log file
	// size, then it is recommended to slightly lower the limit value to ensure that the
	// massive log will not break the limited size, in general, I think this is a small problem,
	// and I don't want to waste time here.
	// In addition, if the log rotation fails, we will continue to write to the original log
	// file, but we will continue to try log rotation until it succeeds.
	if n > 0 && w.max > 0 && atomic.AddUint32(&w.size, uint32(n)) >= w.max {
		w.swap()
	}
	return
}

// Rotate the current log file.
// At the same time trigger the backup log file cleanup (if needed).
func (w *fileWriter) swap() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if atomic.LoadUint32(&w.size) < w.max {
		return
	}
	// We need to make sure that the log file is written to disk before closing the file.
	if err := w.fd.Sync(); err != nil {
		internal.EchoError("Failed to sync log file %s: %s.", w.name, err)
	}
	// On the windows platform, the file must be closed before renaming the file.
	// We need to ignore the error that the file is closed.
	if err := w.fd.Close(); err != nil {
		internal.EchoError("Failed to close log file %s: %s.", w.name, err)
	}
	// Accurate to the nanosecond, it maximizes the assurance that file names are not duplicated.
	suffix := time.Now().Format("20060102150405.000000000")
	if err := os.Rename(w.name, w.name+"."+suffix); err != nil {
		internal.EchoError("Failed to rename log file %s to %s.%s: %s.", w.name, w.name, suffix, err)
	}
	fd, err := os.OpenFile(w.name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		internal.EchoError("Failed to open log file %s: %s.", w.name, err)
		return
	}
	w.fd = fd
	atomic.StoreUint32(&w.size, 0)
	if w.backup > 0 {
		// Try cleaning up old backup log files.
		// If the number of previous backup files is very large, then this may be very
		// time-consuming, so we use asynchronous methods to remove obsolete log backups.
		go w.clearBackups()
	}
}

// Clean up log backup files, if necessary.
// This method ignores any errors generated by cleaning up the backup log file, and all error
// messages will be redirected to standard error.
func (w *fileWriter) clearBackups() {
	matches, err := filepath.Glob(w.name + ".*")
	if err != nil {
		internal.EchoError("Failed to search log backup file %s.*: %s.", w.name, err)
		return
	}
	if n := uint32(len(matches)); n > w.backup {
		// Because we rename log files with timestamps, the oldest log file name is always
		// at the beginning of the slice.
		sort.Strings(matches)
		for i, j := uint32(0), n-w.backup; i < j; i++ {
			if err = os.Remove(matches[i]); err != nil {
				internal.EchoError("Failed to remove log backup file %s: %s.", matches[i], err)
			}
		}
	}
}

// Close is an implementation of the io.WriteCloser interface.
// This method closes the open file descriptor, but keeps the reference pointer to it.
func (w *fileWriter) Close() (err error) {
	w.mu.RLock()
	err = w.fd.Close()
	w.mu.RUnlock()
	return
}
