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
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewFileWriter(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "test.log")

	if w, err := NewFileWriter(name, 1024, 0); err != nil {
		t.Fatal(err)
	} else {
		if w == nil {
			t.Fatal("NewFileWriter(): nil")
		}
		if err = w.Close(); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := NewFileWriter(dir, 1024, 0); err == nil {
		t.Fatal("NewFileWriter(): nil error")
	}
}

func TestNewFileWriter_Error(t *testing.T) {
	dir := t.TempDir()
	// A directory cannot be used as a log file name.
	if _, err := NewFileWriter(dir, 1024, 0); err == nil {
		t.Fatal("NewFileWriter(): nil error")
	}
}

func TestMustNewFileWriter(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "test.log")
	if w := MustNewFileWriter(name, 1024, 0); w == nil {
		t.Fatal("MustNewFileWriter(): nil")
	} else {
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestMustNewFileWriter_Panic(t *testing.T) {
	dir := t.TempDir()
	defer func() {
		if v := recover(); v == nil {
			t.Fatal("MustNewFileWriter(): not panic")
		}
	}()
	// A directory cannot be used as a log file name.
	MustNewFileWriter(dir, 1024, 0)
}

func TestFileWriter(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "test.log")
	w := MustNewFileWriter(name, 1024, 0)

	data := bytes.Repeat([]byte("1"), 1024)
	if n, err := w.Write(data); err != nil {
		t.Fatal(err)
	} else {
		if n != len(data) {
			t.Fatalf("FileWriter.Write(): %d", n)
		}
	}

	if matches, err := filepath.Glob(filepath.Join(dir, "test-*.log")); err != nil {
		t.Fatal(err)
	} else {
		if len(matches) != 1 {
			t.Fatalf("FileWriter.Write(): %v", matches)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestFileWriterWithBackup(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "test.log")
	w := MustNewFileWriter(name, 1000, 2)

	data := bytes.Repeat([]byte("1"), 1024)
	for i := 0; i < 3; i++ {
		if n, err := w.Write(data); err != nil {
			t.Fatal(err)
		} else {
			if n != len(data) {
				t.Fatalf("FileWriter.Write(): %d", n)
			}
		}
	}
	// Delete is asynchronous, we have to wait for system.
	time.Sleep(time.Second * 3)
	if matches, err := filepath.Glob(filepath.Join(dir, "test-*.log")); err != nil {
		t.Fatal(err)
	} else {
		if len(matches) != 2 {
			t.Fatalf("FileWriter.Write(): %v", matches)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestFileWriterWithExistFile(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "test.log")
	if err := os.WriteFile(name, []byte("test"), 0666); err != nil {
		t.Fatal(err)
	}
	w := MustNewFileWriter(name, 1024, 0)

	data := bytes.Repeat([]byte("1"), 1020)
	if n, err := w.Write(data); err != nil {
		t.Fatal(err)
	} else {
		if n != len(data) {
			t.Fatalf("FileWriter.Write(): %d", n)
		}
	}

	if matches, err := filepath.Glob(filepath.Join(dir, "test-*.log")); err != nil {
		t.Fatal(err)
	} else {
		if len(matches) != 1 {
			t.Fatalf("FileWriter.Write(): %v", matches)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestFileWriterWithExistBigFile(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "test.log")
	data := bytes.Repeat([]byte("1"), 2000)
	if err := os.WriteFile(name, data, 0666); err != nil {
		t.Fatal(err)
	}
	w := MustNewFileWriter(name, 1024, 0)

	if matches, err := filepath.Glob(filepath.Join(dir, "test-*.log")); err != nil {
		t.Fatal(err)
	} else {
		if len(matches) != 1 {
			t.Fatalf("FileWriter.Write(): %v", matches)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
}
