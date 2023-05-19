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
	"strings"
	"testing"
)

func testLogCaller(o Log) {
	o.Debug("debug")
	o.Info("info")
}

func TestLogger_EnableCaller(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)
	o.EnableCaller()
	o.Debug("hello") // LINE 34

	got := w.String()
	if !strings.Contains(got, "logger_caller_test.go:34") {
		t.Fatalf("Logger caller: %s", got)
	}

	w.Reset()
	o.EnableCaller(1)
	testLogCaller(o) // LINE 43

	got = w.String()
	if !strings.Contains(got, "logger_caller_test.go:43") {
		t.Fatalf("Logger caller: %s", got)
	}
}

func TestLogger_EnableLevelsCaller(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)
	o.EnableLevelsCaller([]Level{DebugLevel})
	o.Debug("bar") // LINE 57
	o.Info("foo")  // LINE 58

	got := w.String()
	if !strings.Contains(got, "logger_caller_test.go:57") {
		t.Fatalf("Logger caller: %s", got)
	}
	if strings.Contains(got, "logger_caller_test.go:58") {
		t.Fatalf("Logger caller: %s", got)
	}

	w.Reset()
	o.EnableLevelsCaller([]Level{DebugLevel}, 1)
	testLogCaller(o) // LINE 70

	got = w.String()
	if !strings.Contains(got, "logger_caller_test.go:70") {
		t.Fatalf("Logger caller: %s", got)
	}
}

func TestLogger_WithCaller(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	o.WithCaller().Info("test")                            // LINE 84
	o.WithCaller(1).WithCaller().WithCaller().Info("test") // LINE 85

	got := w.String()
	if !strings.Contains(got, "logger_caller_test.go:84") {
		t.Fatalf("Logger caller: %s", got)
	}
	if !strings.Contains(got, "logger_caller_test.go:85") {
		t.Fatalf("Logger caller: %s", got)
	}

	w.Reset()
	testLogCaller(o.WithCaller(1)) // LINE 96
	got = w.String()
	if !strings.Contains(got, "logger_caller_test.go:96") {
		t.Fatalf("Logger caller: %s", got)
	}
}

func TestLogger_InvalidCaller(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)
	o.WithCaller(10).Info("test")

	got := w.String()
	if !strings.Contains(got, "???:0") {
		t.Fatalf("Logger caller: %s", got)
	}
}

func TestLogger_WithCaller_Skip(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)
	o.EnableCaller(1)

	f1 := func() { o.WithCaller(2).Info("test") }
	f2 := func() { f1() }
	f3 := func() { f2() }
	f3() // LINE 126

	got := w.String()
	if !strings.Contains(got, "logger_caller_test.go:126") {
		t.Fatalf("Logger caller: %s", got)
	}
}

func TestLogger_SetCallerSkip(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)
	o.SetCallerSkip(4)

	with := false
	f := func() {
		if with {
			o.WithCaller(1).Debug("debug")
		} else {
			o.Debug("debug")
		}
	}
	f1 := func() { f() }
	f2 := func() { f1() }
	f3 := func() { f2() }
	f4 := func() { f3() }
	f4() // LINE 153

	if got := w.String(); strings.Contains(got, "logger_caller_test.go:153") {
		t.Fatalf("Logger caller: %s", got)
	}

	with = true
	w.Reset()
	f4() // LINE 161
	if got := w.String(); !strings.Contains(got, "logger_caller_test.go:161") {
		t.Fatalf("Logger caller: %s", got)
	}
}

func TestLogger_SetLongCaller(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)
	o.SetFormatter(FormatterFunc(func(e Entity, b *bytes.Buffer) error {
		b.WriteString(e.Caller())
		return nil
	}))

	f1 := func() { o.WithCaller(4).Debug("debug") }
	f2 := func() { f1() }
	f3 := func() { f2() }
	f4 := func() { f3() }
	f4() // LINE 181

	if got := w.String(); !strings.HasSuffix(got, "logger_caller_test.go:181") {
		t.Fatalf("Logger caller: %s", got)
	}
}
