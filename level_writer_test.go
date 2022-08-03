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
	"testing"
)

func TestNewLevelWriter(t *testing.T) {
	o := New("test")
	if NewLevelWriter(InfoLevel, o.AsLog()) == nil {
		t.Fatal("NewLevelWriter(): nil")
	}
}

func TestLevelWriter_Write(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)
	o.SetFormatter(FormatterFunc(func(e Entity, b *bytes.Buffer) error {
		b.WriteString(e.Level().String() + " " + e.Message())
		return nil
	}))
	lw := NewLevelWriter(DebugLevel, o.AsLog())

	items := [][]byte{[]byte("test"), []byte("test\n")}
	for i, j := 0, len(items); i < j; i++ {
		w.Reset()
		if n, err := lw.Write(items[i]); err != nil {
			t.Fatalf("LevelWriter.Write(): %s", err)
		} else {
			if n != len(items[i]) {
				t.Fatalf("LevelWriter.Write(): %d", n)
			}
		}
		if want, got := "debug test", w.String(); want != got {
			t.Fatalf("LevelWriter.Write(): got %q, want %q", got, want)
		}
	}
}
