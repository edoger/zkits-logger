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
	"encoding/json"
	"os"
	"testing"
	"time"
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

func TestLogger_Log(t *testing.T) {
	w := new(bytes.Buffer)
	g := New("test")
	g.SetOutput(w)

	type O struct {
		Level   string                 `json:"level"`
		Logger  string                 `json:"logger"`
		Message string                 `json:"message"`
		Time    string                 `json:"time"`
		Fields  map[string]interface{} `json:"fields"`
	}

	var o *O

	check := func(o *O, message string, level Level) {
		if err := json.Unmarshal(w.Bytes(), o); err != nil {
			t.Fatal(err)
		}
		if o.Level != level.String() {
			t.Fatal(o.Level)
		}
		if o.Logger != "test" {
			t.Fatal(o.Logger)
		}
		if o.Message != message {
			t.Fatal(o.Message, message)
		}
		if _, err := time.Parse("2006-01-02 15:04:05", o.Time); err != nil {
			t.Fatal(err)
		}
		if o.Fields == nil || len(o.Fields) != 0 {
			t.Fatal(o.Fields)
		}
	}
	do := func(f func()) {
		w.Reset()
		o = new(O)
		f()
	}

	// -------------- TraceLevel -----------------
	g.SetLevel(TraceLevel)

	do(func() {
		g.Trace("foo")
		check(o, "foo", TraceLevel)
	})

	do(func() {
		g.Traceln("bar")
		check(o, "bar", TraceLevel)
	})

	do(func() {
		g.Tracef("bar-%d", 1)
		check(o, "bar-1", TraceLevel)
	})

	do(func() {
		g.Print("foo")
		check(o, "foo", TraceLevel)
	})

	do(func() {
		g.Println("bar")
		check(o, "bar", TraceLevel)
	})

	do(func() {
		g.Printf("bar-%d", 1)
		check(o, "bar-1", TraceLevel)
	})

	// -------------- DebugLevel -----------------
	g.SetLevel(DebugLevel)

	do(func() {
		g.Trace("test-1")
		g.Traceln("test-2")
		g.Tracef("test-%d", 3)
		g.Print("test-4")
		g.Println("test-5")
		g.Printf("test-%d", 6)

		if w.Len() != 0 {
			t.Fatal(w.Len())
		}
	})

}
