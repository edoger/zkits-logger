// Copyright 2020 The ZKits Project Authors.
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
	"testing"
)

func TestNew(t *testing.T) {
	g := New("test")

	if g == nil {
		t.Fatal("New() return nil")
	}
}

func TestLogger_Name(t *testing.T) {
	g := New("test")

	if name := g.Name(); name != "test" {
		t.Fatal(name)
	}
}

func TestLogger_Level(t *testing.T) {
	g := New("test")

	if level := g.GetLevel(); level != InfoLevel {
		t.Fatal(level.String())
	}

	if g.SetLevel(Level(0)) == nil {
		t.Fatal("Logger.SetLevel() return nil")
	}
	if level := g.GetLevel(); level != InfoLevel {
		t.Fatal(level.String())
	}

	if g.SetLevel(DebugLevel) == nil {
		t.Fatal("Logger.SetLevel() return nil")
	}
	if level := g.GetLevel(); level != DebugLevel {
		t.Fatal(level.String())
	}
}

func TestLogger_Output(t *testing.T) {
	g := New("test")

	if w := g.GetOutput(); w == nil {
		t.Fatal("Logger.GetOutput() return nil")
	} else {
		f, ok := w.(*os.File)
		if !ok {
			t.Fatal("Default output not is *os.File")
		}
		if f.Name() != os.Stdout.Name() {
			t.Fatal(f.Name())
		}
	}

	if g.SetOutput(bytes.NewBufferString("foo")) == nil {
		t.Fatal("Logger.SetOutput() return nil")
	}

	if w := g.GetOutput(); w == nil {
		t.Fatal("Logger.GetOutput() return nil")
	} else {
		b, ok := w.(*bytes.Buffer)
		if !ok {
			t.Fatal("Default output not is *bytes.Buffer")
		}
		if b.String() != "foo" {
			t.Fatal(b.String())
		}
	}
}
