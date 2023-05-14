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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/edoger/zkits-logger/internal"
)

const dirPerm os.FileMode = 0766
const filePerm os.FileMode = 0666
const fileFlag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
const backupTimeFormat = "2006-01-02T15:04:05.000"

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
	w := &fileWriter{path: name, max: max, backup: backup, clear: make(chan struct{}, 1)}
	if err := w.open(); err != nil {
		return nil, err
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
	mu     sync.Mutex
	file   *os.File
	path   string
	size   uint32 // The current log file size.
	max    uint32 // The maximum size of the log file.
	backup uint32 // The maximum number of backup log files.
	once   sync.Once
	clear  chan struct{}
}

// Write is an implementation of the io.WriteCloser interface, used to write a single
// log data to a file.
func (w *fileWriter) Write(b []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file == nil {
		err = w.open()
		if err != nil {
			return
		}
	}
	n, err = w.file.Write(b)
	if w.max > 0 {
		w.size += uint32(n)
		if w.size >= w.max {
			w.rotate()
		}
	}
	return
}

// Close is an implementation of the io.WriteCloser interface.
func (w *fileWriter) Close() (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		defer func() { w.file, w.size = nil, 0 }()
		err = w.file.Sync()
		if err2 := w.file.Close(); err == nil {
			err = err2
		}
	}
	return
}

func (w *fileWriter) open() error {
	dir, name, ext := splitFilePath(w.path)
	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		return err
	}
	var info os.FileInfo
	var file *os.File
	if info, err = os.Stat(w.path); err != nil {
		if os.IsNotExist(err) {
			file, err = os.OpenFile(w.path, fileFlag, filePerm)
			if err != nil {
				return err
			}
			w.file, w.size = file, 0
			return nil
		}
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("path %s exists, but not is regular file", w.path)
	}
	if w.max > 0 && uint32(info.Size()) > w.max {
		if err = os.Rename(w.path, newBackupFileName(dir, name, ext)); err != nil {
			return err
		}
		w.clean()
	}
	if file, err = os.OpenFile(w.path, fileFlag, filePerm); err != nil {
		return err
	}
	if info, err = file.Stat(); err != nil {
		_ = file.Close()
		return err
	}
	w.file = file
	w.size = uint32(info.Size())
	return nil
}

func (w *fileWriter) clean() {
	w.once.Do(w.sweeper)
	select {
	case w.clear <- struct{}{}:
	default:
	}
}

func (w *fileWriter) rotate() {
	if err := w.file.Sync(); err != nil {
		internal.EchoError("Failed to sync %s: %s.", w.path, err)
	}
	if err := w.file.Close(); err != nil {
		internal.EchoError("Failed to close %s: %s.", w.path, err)
	}
	w.file, w.size = nil, 0
	dir, name, ext := splitFilePath(w.path)
	if err := os.Rename(w.path, newBackupFileName(dir, name, ext)); err != nil {
		internal.EchoError("Failed to rename %s: %s.", w.path, err)
	}
	w.clean()
}

func (w *fileWriter) sweeper() {
	go func() {
		for {
			<-w.clear
			dir, name, ext := splitFilePath(w.path)
			if items, err := os.ReadDir(dir); err == nil {
				if len(items) == 0 {
					continue
				}
				base, files := filepath.Base(w.path), make([]string, 0)
				for i, j := 0, len(items); i < j; i++ {
					if items[i].Type().IsRegular() && isBackupFileName(items[i].Name(), base, name+"-", ext) {
						files = append(files, filepath.Join(dir, items[i].Name()))
					}
				}
				if n := uint32(len(files)); n > w.backup {
					removeFiles(files[:n-w.backup])
				}
			} else {
				internal.EchoError("Call os.ReadDir() with dir %s failed: %s.", dir, err)
			}
		}
	}()
}

// Delete the given list of files.
func removeFiles(files []string) {
	for i, j := 0, len(files); i < j; i++ {
		if err := os.Remove(files[i]); err != nil && !os.IsNotExist(err) {
			internal.EchoError("Failed to remove %s: %s.", files[i], err)
		}
	}
}

// Split the basic metadata for the given path.
func splitFilePath(path string) (dir, name, ext string) {
	base := filepath.Base(path)
	ext = filepath.Ext(base)
	dir, name = filepath.Dir(path), strings.TrimSuffix(base, ext)
	return
}

// Create a new log backup file name.
func newBackupFileName(dir, name, ext string) string {
	return filepath.Join(dir, name+"-"+time.Now().Local().Format(backupTimeFormat)+ext)
}

// Determines if the given filename is a log backup file.
func isBackupFileName(target, base, name, ext string) bool {
	if target == base || !strings.HasPrefix(target, name) {
		return false
	}
	target = target[len(name):]
	if !strings.HasSuffix(target, ext) {
		return false
	}
	target = target[:len(target)-len(ext)]
	if len(target) != len(backupTimeFormat) {
		return false
	}
	if _, err := time.Parse(backupTimeFormat, target); err != nil {
		return false
	}
	return true
}
