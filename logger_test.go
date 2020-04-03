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

	do(func() {
		g.Debug("foo")
		check(o, "foo", DebugLevel)
	})

	do(func() {
		g.Debugln("bar")
		check(o, "bar", DebugLevel)
	})

	do(func() {
		g.Debugf("bar-%d", 1)
		check(o, "bar-1", DebugLevel)
	})

	// -------------- InfoLevel -----------------
	g.SetLevel(InfoLevel)

	do(func() {
		g.Debug("test-1")
		g.Debugln("test-2")
		g.Debugf("test-%d", 3)

		if w.Len() != 0 {
			t.Fatal(w.Len())
		}
	})

	do(func() {
		g.Info("foo")
		check(o, "foo", InfoLevel)
	})

	do(func() {
		g.Infoln("bar")
		check(o, "bar", InfoLevel)
	})

	do(func() {
		g.Infof("bar-%d", 1)
		check(o, "bar-1", InfoLevel)
	})

	// -------------- WarnLevel -----------------
	g.SetLevel(WarnLevel)

	do(func() {
		g.Info("test-1")
		g.Infoln("test-2")
		g.Infof("test-%d", 3)

		if w.Len() != 0 {
			t.Fatal(w.Len())
		}
	})

	do(func() {
		g.Warn("foo")
		check(o, "foo", WarnLevel)
	})

	do(func() {
		g.Warnln("bar")
		check(o, "bar", WarnLevel)
	})

	do(func() {
		g.Warnf("bar-%d", 1)
		check(o, "bar-1", WarnLevel)
	})

	do(func() {
		g.Warning("foo")
		check(o, "foo", WarnLevel)
	})

	do(func() {
		g.Warningln("bar")
		check(o, "bar", WarnLevel)
	})

	do(func() {
		g.Warningf("bar-%d", 1)
		check(o, "bar-1", WarnLevel)
	})

	// -------------- ErrorLevel -----------------
	g.SetLevel(ErrorLevel)

	do(func() {
		g.Warn("test-1")
		g.Warnln("test-2")
		g.Warnf("test-%d", 3)
		g.Warning("test-4")
		g.Warningln("test-5")
		g.Warningf("test-%d", 6)

		if w.Len() != 0 {
			t.Fatal(w.Len())
		}
	})

	do(func() {
		g.Error("foo")
		check(o, "foo", ErrorLevel)
	})

	do(func() {
		g.Errorln("bar")
		check(o, "bar", ErrorLevel)
	})

	do(func() {
		g.Errorf("bar-%d", 1)
		check(o, "bar-1", ErrorLevel)
	})

	// -------------- FatalLevel -----------------
	g.SetLevel(FatalLevel)

	var exitCode int

	g.WithExitFunc(func(i int) { exitCode = i })

	checkExit := func(code int) {
		if exitCode != code {
			t.Fatalf("No exit: %d", exitCode)
		}
	}

	do(func() {
		g.Error("test-1")
		g.Errorln("test-2")
		g.Errorf("test-%d", 3)

		if w.Len() != 0 {
			t.Fatal(w.Len())
		}
	})

	do(func() {
		exitCode = 0
		g.Fatal("foo")
		check(o, "foo", FatalLevel)
		checkExit(1)
	})

	do(func() {
		exitCode = 0
		g.Fatalln("bar")
		check(o, "bar", FatalLevel)
		checkExit(1)
	})

	do(func() {
		exitCode = 0
		g.Fatalf("bar-%d", 1)
		check(o, "bar-1", FatalLevel)
		checkExit(1)
	})

	// -------------- PanicLevel -----------------
	g.SetLevel(PanicLevel)

	var panicValue interface{}

	checkPanic := func() {
		if panicValue == nil {
			t.Fatal("No panic")
		}
		if _, ok := panicValue.(Summary); !ok {
			t.Fatalf("Panic value %T", panicValue)
		}
	}

	call := func(f func()) {
		defer func() { panicValue = recover() }()
		f()
	}

	do(func() {
		exitCode = 0
		g.Fatal("test-1")
		g.Fatalln("test-2")
		g.Fatalf("test-%d", 3)

		if w.Len() != 0 {
			t.Fatal(w.Len())
		}
		checkExit(0)
	})

	do(func() {
		panicValue = nil
		call(func() { g.Panic("foo") })
		check(o, "foo", PanicLevel)
		checkPanic()
	})

	do(func() {
		panicValue = nil
		call(func() { g.Panicln("bar") })
		check(o, "bar", PanicLevel)
		checkPanic()
	})

	do(func() {
		panicValue = nil
		call(func() { g.Panicf("bar-%d", 1) })
		check(o, "bar-1", PanicLevel)
		checkPanic()
	})
}
