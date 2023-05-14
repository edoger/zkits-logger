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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/edoger/zkits-logger/internal"
)

func TestNew(t *testing.T) {
	o := New("test")

	if o == nil {
		t.Fatal("New(): nil")
	}
}

func TestLogger_Name(t *testing.T) {
	o := New("test")

	if name := o.Name(); name != "test" {
		t.Fatalf("Logger.Name(): %s", name)
	}
}

func TestLogger_Level(t *testing.T) {
	o := New("test")

	if level := o.GetLevel(); level != TraceLevel {
		t.Fatalf("Logger.GetLevel(): %s", level.String())
	}
	if o.SetLevel(DebugLevel) == nil {
		t.Fatal("Logger.SetLevel() return nil")
	}
	if level := o.GetLevel(); level != DebugLevel {
		t.Fatalf("Logger.GetLevel(): %s", level.String())
	}
}

func TestLogger_SetOutput(t *testing.T) {
	o := New("test")
	w := new(bytes.Buffer)

	if o.SetOutput(w) == nil {
		t.Fatal("Logger.SetOutput(): nil")
	}
	if o.SetOutput(nil) == nil {
		t.Fatal("Logger.SetOutput(nil): nil")
	}
}

func TestLogger_SetLevelOutput(t *testing.T) {
	o := New("test")
	w := new(bytes.Buffer)

	if o.SetLevelOutput(InfoLevel, w) == nil {
		t.Fatal("Logger.SetLevelOutput(): nil")
	}
	if o.SetLevelOutput(InfoLevel, nil) == nil {
		t.Fatal("Logger.SetLevelOutput(Level, nil): nil")
	}
}

func TestLogger_SetLevelsOutput(t *testing.T) {
	o := New("test")
	w := new(bytes.Buffer)

	if o.SetLevelsOutput([]Level{InfoLevel, WarnLevel}, w) == nil {
		t.Fatal("Logger.SetLevelsOutput(): nil")
	}
	if o.SetLevelsOutput([]Level{InfoLevel, WarnLevel}, nil) == nil {
		t.Fatal("Logger.SetLevelsOutput(Level, nil): nil")
	}
}

func TestLogger_SetNowFunc(t *testing.T) {
	o := New("test")
	f := func() time.Time { return time.Now() }

	if o.SetNowFunc(f) == nil {
		t.Fatal("Logger.SetNowFunc(): nil")
	}
	if o.SetNowFunc(nil) == nil {
		t.Fatal("Logger.SetNowFunc(nil): nil")
	}
}

func TestLogger_SetExitFunc(t *testing.T) {
	o := New("test")
	f := func(int) {}

	if o.SetExitFunc(f) == nil {
		t.Fatal("Logger.SetExitFunc(): nil")
	}
	if o.SetExitFunc(nil) == nil {
		t.Fatal("Logger.SetExitFunc(nil): nil")
	}
}

func TestLogger_SetPanicFunc(t *testing.T) {
	o := New("test")
	f := func(string) {}

	if o.SetPanicFunc(f) == nil {
		t.Fatal("Logger.SetPanicFunc(): nil")
	}
	if o.SetPanicFunc(nil) == nil {
		t.Fatal("Logger.SetPanicFunc(nil): nil")
	}
}

func TestLogger_SetFormatter(t *testing.T) {
	o := New("test")
	f := FormatterFunc(func(Entity, *bytes.Buffer) error { return nil })

	if o.SetFormatter(f) == nil {
		t.Fatal("Logger.SetFormatter(): nil")
	}
	if o.SetFormatter(nil) == nil {
		t.Fatal("Logger.SetFormatter(nil): nil")
	}
}

func TestLogger_SetDefaultTimeFormat(t *testing.T) {
	o := New("test")
	if o.SetDefaultTimeFormat("2006-01-02 15:04:05") == nil {
		t.Fatal("Logger.SetDefaultTimeFormat(): nil")
	}
	if o.SetDefaultTimeFormat("") == nil {
		t.Fatal("Logger.SetDefaultTimeFormat(): nil")
	}
}

func TestLogger_Log(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")

	o.SetOutput(w)
	o.SetExitFunc(nil)  // Disable
	o.SetPanicFunc(nil) // Disable
	o.SetLevel(TraceLevel)

	now := time.Now()
	o.SetNowFunc(func() time.Time { return now })

	do := func(s string, fs ...func(Logger) (Level, string)) {
		buf := new(bytes.Buffer)
		for _, f := range fs {
			w.Reset()
			buf.Reset()
			level, message := f(o)
			if level.IsValid() {
				err := json.NewEncoder(buf).Encode(map[string]interface{}{
					"name":    "test",
					"time":    now.Format(time.RFC3339),
					"level":   level.String(),
					"message": message,
				})
				if err != nil {
					t.Fatalf("%s: %s", s, err)
				}
				if !bytes.Equal(w.Bytes(), buf.Bytes()) {
					t.Fatalf("%s: %s -- %s", s, w.String(), buf.String())
				}
			} else {
				// No log
				if got := w.String(); got != "" {
					t.Fatalf("%s: %s", s, got)
				}
			}
		}
	}

	// -------------- TraceLevel -----------------

	do("TraceLevel", func(o Logger) (Level, string) {
		o.Trace("foo")
		return TraceLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Traceln("foo")
		return TraceLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Tracef("foo-%s", "bar")
		return TraceLevel, fmt.Sprintf("foo-%s", "bar")
	}, func(o Logger) (Level, string) {
		o.Print("foo")
		return TraceLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Println("foo")
		return TraceLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Printf("foo-%s", "bar")
		return TraceLevel, fmt.Sprintf("foo-%s", "bar")
	})

	// -------------- DebugLevel -----------------

	do("DebugLevel", func(o Logger) (Level, string) {
		o.Debug("foo")
		return DebugLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Debugln("foo")
		return DebugLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Debugf("foo-%s", "bar")
		return DebugLevel, fmt.Sprintf("foo-%s", "bar")
	})

	// -------------- InfoLevel -----------------

	do("InfoLevel", func(o Logger) (Level, string) {
		o.Info("foo")
		return InfoLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Infoln("foo")
		return InfoLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Infof("foo-%s", "bar")
		return InfoLevel, fmt.Sprintf("foo-%s", "bar")
	}, func(o Logger) (Level, string) {
		o.Echo("foo")
		return InfoLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Echoln("foo")
		return InfoLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Echof("foo-%s", "bar")
		return InfoLevel, fmt.Sprintf("foo-%s", "bar")
	})

	// -------------- WarnLevel -----------------

	do("WarnLevel", func(o Logger) (Level, string) {
		o.Warn("foo")
		return WarnLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Warnln("foo")
		return WarnLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Warnf("foo-%s", "bar")
		return WarnLevel, fmt.Sprintf("foo-%s", "bar")
	}, func(o Logger) (Level, string) {
		o.Warning("foo")
		return WarnLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Warningln("foo")
		return WarnLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Warningf("foo-%s", "bar")
		return WarnLevel, fmt.Sprintf("foo-%s", "bar")
	})

	// -------------- ErrorLevel -----------------

	do("ErrorLevel", func(o Logger) (Level, string) {
		o.Error("foo")
		return ErrorLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Errorln("foo")
		return ErrorLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Errorf("foo-%s", "bar")
		return ErrorLevel, fmt.Sprintf("foo-%s", "bar")
	})

	// -------------- FatalLevel -----------------

	do("FatalLevel", func(o Logger) (Level, string) {
		o.Fatal("foo")
		return FatalLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Fatalln("foo")
		return FatalLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Fatalf("foo-%s", "bar")
		return FatalLevel, fmt.Sprintf("foo-%s", "bar")
	})

	// -------------- PanicLevel -----------------

	do("PanicLevel", func(o Logger) (Level, string) {
		o.Panic("foo")
		return PanicLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Panicln("foo")
		return PanicLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Panicf("foo-%s", "bar")
		return PanicLevel, fmt.Sprintf("foo-%s", "bar")
	})

	// Test a higher level of log.
	o.SetLevel(ErrorLevel)

	do("Use ErrorLevel", func(o Logger) (Level, string) {
		o.Trace("foo")
		return 0, "" // No log
	}, func(o Logger) (Level, string) {
		o.Debug("foo")
		return 0, "" // No log
	}, func(o Logger) (Level, string) {
		o.Info("foo")
		return 0, "" // No log
	}, func(o Logger) (Level, string) {
		o.Warn("foo")
		return 0, "" // No log
	}, func(o Logger) (Level, string) {
		o.Error("foo")
		return ErrorLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Fatal("foo")
		return FatalLevel, "foo"
	}, func(o Logger) (Level, string) {
		o.Panic("foo")
		return PanicLevel, "foo"
	})
}

func TestLogger_LogPanic(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	defer func() {
		if recover() == nil {
			t.Fatal("No panic")
		}
	}()

	o.Panic("foo")
}

func TestLoggerWithLevel(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(ErrorLevel)

	o.Log(InfoLevel, "foo")
	o.Logf(InfoLevel, "%d", 1)
	o.Logln(InfoLevel, "bar")

	if got := w.String(); got != "" {
		t.Fatalf("Logger: %s", got)
	}
}

func TestLogger_Hook(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	var message string

	o.AddHookFunc([]Level{TraceLevel}, func(s Summary) error {
		message = s.Message()
		return nil
	})

	o.Trace("foo")

	if message != "foo" {
		t.Fatalf("Log hook not exec: %s", message)
	}
}

func TestLogger_WithContext(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	var ctx context.Context

	o.AddHookFunc([]Level{TraceLevel}, func(s Summary) error {
		ctx = s.Context()
		return nil
	})

	o.Trace("foo") // Without context
	if ctx == nil {
		t.Fatal("Context: nil")
	}
	if got := ctx.Value("key"); got != nil {
		t.Fatalf("Context: %v", got)
	}

	key := struct{}{}
	o.WithContext(context.Background()).WithContext(context.WithValue(context.Background(), key, "foo")).Trace("foo")
	if ctx == nil {
		t.Fatal("Context: nil")
	}
	if got := ctx.Value(key).(string); got != "foo" {
		t.Fatalf("Context: %s", got)
	}
}

func TestLogger_WithFields(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	var fields map[string]interface{}

	o.AddHookFunc([]Level{TraceLevel}, func(s Summary) error {
		fields = s.Fields()
		return nil
	})

	o.Trace("foo") // Without fields
	if len(fields) != 0 {
		t.Fatalf("Fields: %v", fields)
	}

	o.WithFields(map[string]interface{}{"key": "bar"}).WithFields(map[string]interface{}{"key": "foo"}).Trace("foo")
	if len(fields) != 1 {
		t.Fatalf("Fields: %v", fields)
	}
	if got := fields["key"].(string); got != "foo" {
		t.Fatalf("Fields: %s", got)
	}
}

func TestLogger_WithField(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	var fields map[string]interface{}

	o.AddHookFunc([]Level{TraceLevel}, func(s Summary) error {
		fields = s.Fields()
		return nil
	})

	o.Trace("foo") // Without field
	if len(fields) != 0 {
		t.Fatalf("Fields: %v", fields)
	}

	o.WithField("key", "foo").Trace("foo")
	if len(fields) != 1 {
		t.Fatalf("Fields: %v", fields)
	}
	if got := fields["key"].(string); got != "foo" {
		t.Fatalf("Fields: %s", got)
	}
}

func TestLogger_WithMessagePrefix(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	o.SetFormatter(FormatterFunc(func(e Entity, b *bytes.Buffer) error {
		b.WriteString(e.Message())
		return nil
	}))

	o.WithMessagePrefix("Prefix: ").Trace("foo")
	want := "Prefix: foo"
	if got := w.String(); got != want {
		t.Fatalf("WithMessagePrefix: got %q, want %q", got, want)
	}
}

func TestLogger_WithError(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	var fields map[string]interface{}

	o.AddHookFunc([]Level{TraceLevel}, func(s Summary) error {
		fields = s.Fields()
		return nil
	})

	o.Trace("foo") // Without error
	if len(fields) != 0 {
		t.Fatalf("Fields: %v", fields)
	}

	o.WithError(errors.New("error")).Trace("foo")
	if len(fields) != 1 {
		t.Fatalf("Fields: %v", fields)
	}
	if got := fields["error"].(error); got.Error() != "error" {
		t.Fatalf("Fields: %s", got)
	}
}

func TestLogger_Formatter(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	o.SetFormatter(FormatterFunc(func(e Entity, b *bytes.Buffer) error {
		b.WriteString("formatter")
		return nil
	}))

	o.Trace("foo")
	if got := w.String(); got != "formatter" {
		t.Fatalf("Formatter: %s", got)
	}
}

func TestLoggerFormatterError(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	buf := new(bytes.Buffer)

	internal.ErrorWriter = buf
	defer func() { internal.ErrorWriter = os.Stderr }()

	o.SetFormatter(FormatterFunc(func(Entity, *bytes.Buffer) error {
		return errors.New("formatter")
	}))

	o.Trace("foo")

	if got := buf.String(); got != "(test) Failed to format log: formatter\n" {
		t.Fatalf("Formatter: %s", got)
	}
}

func TestLoggerHookError(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	buf := new(bytes.Buffer)

	internal.ErrorWriter = buf
	defer func() { internal.ErrorWriter = os.Stderr }()

	o.AddHookFunc([]Level{TraceLevel}, func(Summary) error {
		return errors.New("hook")
	})

	o.Trace("foo")

	if got := buf.String(); got != "(test) Failed to fire log hook: hook\n" {
		t.Fatalf("Hook: %s", got)
	}
}

type testErrorWriter string

func (s testErrorWriter) Write([]byte) (int, error) {
	return 0, errors.New(string(s))
}

func TestLoggerWriterError(t *testing.T) {
	w := testErrorWriter("writer")
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	buf := new(bytes.Buffer)

	internal.ErrorWriter = buf
	defer func() { internal.ErrorWriter = os.Stderr }()

	o.Trace("foo")

	if got := buf.String(); got != "(test) Failed to write log: writer\n" {
		t.Fatalf("Writer: %s", got)
	}
}

func TestLoggerHookUseHookBag(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	var ok bool

	bag := NewHookBag()
	bag.Add(NewHookFromFunc([]Level{TraceLevel}, func(Summary) error {
		ok = true
		return nil
	}))

	o.AddHook(bag)
	o.Trace("foo")

	if !ok {
		t.Fatalf("UseHookBag: %v", ok)
	}
}

func TestLogger_LevelOutput(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	w2 := new(bytes.Buffer)
	o.SetLevelOutput(TraceLevel, w2)

	o.Trace("foo")
	if w2.Len() == 0 {
		t.Fatal("LevelOutput: empty output")
	}
}

func TestLogger_BigLog(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	var builder strings.Builder
	for i := 0; i < 5000; i++ {
		builder.WriteByte('x')
	}
	o.Print(builder.String()) // Big log
	if w.Len() < 5000 {
		t.Fatalf("Big log: %d", w.Len())
	}
	// For logs exceeding 4KB, the buffer will not be reused.
	n := o.(*logger).core.pool.Get().(*logEntity).buffer.Cap()
	if n > 4096 {
		t.Fatalf("Big log: %d", n)
	}
}

func TestLogger_Interceptor(t *testing.T) {
	w1 := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w1)
	o.SetLevel(TraceLevel)

	w2 := new(bytes.Buffer)
	interceptor := func(summary Summary, writer io.Writer) (int, error) {
		return w2.Write([]byte(summary.Message())) // message only
	}
	if o.SetOutputInterceptor(interceptor) == nil {
		t.Fatal("Logger.SetOutputInterceptor(): nil")
	}

	o.Echo("foo")
	if got := w1.String(); got != "" {
		t.Fatalf("Logger.SetOutputInterceptor(): got %s", got)
	}
	if got := w2.String(); got != "foo" {
		t.Fatalf("Logger.SetOutputInterceptor(): got %s", got)
	}
}

func TestLogger_AsLog(t *testing.T) {
	o := New("test")
	l := o.AsLog()
	if l == nil {
		t.Fatal("Logger.AsLog(): nil")
	}
	if _, ok := l.(Logger); ok {
		t.Fatal("Logger.AsLog(): Got Logger instance.")
	}
}

func TestLogger_AsStandardLogger(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)
	o.SetFormatter(FormatterFunc(func(e Entity, b *bytes.Buffer) error {
		b.WriteString(e.Level().String() + " " + e.Message())
		return nil
	}))

	l := o.AsStandardLogger()
	if l == nil {
		t.Fatal("Logger.AsStandardLogger(): nil")
	}
	l.Print("test")
	want := "info test"
	if got := w.String(); got != want {
		t.Fatalf("Logger.AsStandardLogger(): got %q, want %q", got, want)
	}
}

func TestLogger_IsLevelEnabled(t *testing.T) {
	o := New("test")
	o.SetLevel(TraceLevel)

	items := []struct {
		Given Level
		Want  bool
	}{
		{TraceLevel, true},
		{DebugLevel, true},
		{InfoLevel, true},
		{WarnLevel, true},
		{ErrorLevel, true},
		{FatalLevel, true},
		{PanicLevel, true},
		{Level(0), false},
		{Level(10000), false},
	}
	for _, item := range items {
		if got := o.IsLevelEnabled(item.Given); got != item.Want {
			t.Fatalf("Logger.IsLevelEnabled(): got %v, want %v", got, item.Want)
		}
	}
	if !o.IsPanicLevelEnabled() {
		t.Fatal("Logger.IsPanicLevelEnabled(): got false, want true")
	}
	if !o.IsFatalLevelEnabled() {
		t.Fatal("Logger.IsFatalLevelEnabled(): got false, want true")
	}
	if !o.IsErrorLevelEnabled() {
		t.Fatal("Logger.IsErrorLevelEnabled(): got false, want true")
	}
	if !o.IsWarnLevelEnabled() {
		t.Fatal("Logger.IsWarnLevelEnabled(): got false, want true")
	}
	if !o.IsInfoLevelEnabled() {
		t.Fatal("Logger.IsInfoLevelEnabled(): got false, want true")
	}
	if !o.IsDebugLevelEnabled() {
		t.Fatal("Logger.IsDebugLevelEnabled(): got false, want true")
	}
	if !o.IsTraceLevelEnabled() {
		t.Fatal("Logger.IsTraceLevelEnabled(): got false, want true")
	}

	o.SetLevel(InfoLevel)

	items = []struct {
		Given Level
		Want  bool
	}{
		{TraceLevel, false},
		{DebugLevel, false},
		{InfoLevel, true},
		{WarnLevel, true},
		{ErrorLevel, true},
		{FatalLevel, true},
		{PanicLevel, true},
		{Level(0), false},
		{Level(10000), false},
	}
	for _, item := range items {
		if got := o.IsLevelEnabled(item.Given); got != item.Want {
			t.Fatalf("Logger.IsLevelEnabled(): got %v, want %v", got, item.Want)
		}
	}
	if !o.IsPanicLevelEnabled() {
		t.Fatal("Logger.IsPanicLevelEnabled(): got false, want true")
	}
	if !o.IsFatalLevelEnabled() {
		t.Fatal("Logger.IsFatalLevelEnabled(): got false, want true")
	}
	if !o.IsErrorLevelEnabled() {
		t.Fatal("Logger.IsErrorLevelEnabled(): got false, want true")
	}
	if !o.IsWarnLevelEnabled() {
		t.Fatal("Logger.IsWarnLevelEnabled(): got false, want true")
	}
	if !o.IsInfoLevelEnabled() {
		t.Fatal("Logger.IsInfoLevelEnabled(): got false, want true")
	}
	if o.IsDebugLevelEnabled() {
		t.Fatal("Logger.IsDebugLevelEnabled(): got true, want false")
	}
	if o.IsTraceLevelEnabled() {
		t.Fatal("Logger.IsTraceLevelEnabled(): got true, want false")
	}
}

func TestLogger_SetLevelWithInvalidLevel(t *testing.T) {
	o := New("test")

	if level := o.GetLevel(); level != TraceLevel {
		t.Fatalf("Logger.GetLevel(): %s", level.String())
	}
	invalidLevel := Level(10000)
	if invalidLevel.IsValid() {
		t.Fatal("Level.IsValid() return true")
	}
	if o.SetLevel(invalidLevel) == nil {
		t.Fatal("Logger.SetLevel() return nil")
	}
	if level := o.GetLevel(); level != TraceLevel {
		t.Fatalf("Logger.GetLevel(): %s", level.String())
	}
}

func TestLogger_SetLevelString(t *testing.T) {
	o := New("test")

	if level := o.GetLevel(); level != TraceLevel {
		t.Fatalf("Logger.GetLevel(): %s", level.String())
	}

	if err := o.SetLevelString("info"); err != nil {
		t.Fatalf("Logger.SetLevelString(): %s", err)
	}
	if level := o.GetLevel(); level != InfoLevel {
		t.Fatalf("Logger.GetLevel(): %s", level.String())
	}

	if err := o.SetLevelString("invalid-level-string"); err == nil {
		t.Fatal("Logger.SetLevelString(): no error")
	}
	if level := o.GetLevel(); level != InfoLevel {
		t.Fatalf("Logger.GetLevel(): %s", level.String())
	}
}

func TestLogger_ForceSetLevelString(t *testing.T) {
	o := New("test")

	if level := o.GetLevel(); level != TraceLevel {
		t.Fatalf("Logger.GetLevel(): %s", level.String())
	}

	if o.ForceSetLevelString("info") == nil {
		t.Fatal("Logger.ForceSetLevelString(): return nil")
	}
	if level := o.GetLevel(); level != InfoLevel {
		t.Fatalf("Logger.GetLevel(): %s", level.String())
	}

	if o.ForceSetLevelString("invalid-level-string") == nil {
		t.Fatal("Logger.ForceSetLevelString(): return nil")
	}
	if level := o.GetLevel(); level != InfoLevel {
		t.Fatalf("Logger.GetLevel(): %s", level.String())
	}
}

func TestLogger_WithFieldPairs(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	var fields map[string]interface{}

	o.AddHookFunc([]Level{TraceLevel}, func(s Summary) error {
		fields = s.Fields()
		return nil
	})

	o.Trace("foo") // Without pairs
	if len(fields) != 0 {
		t.Fatalf("Fields: %v", fields)
	}

	o.WithFieldPairs("foo", "foo").Trace("foo")
	if len(fields) != 1 {
		t.Fatal("Fields: nil")
	}
	if got := fields["foo"].(string); got != "foo" {
		t.Fatalf("Fields: %s", got)
	}

	o.WithFieldPairs("foo", "foo", "bar", "bar", "foo", "test").Trace("foo")
	if len(fields) != 2 {
		t.Fatalf("Fields: %v", fields)
	}
	if got := fields["foo"].(string); got != "test" {
		t.Fatalf("Fields: %s", got)
	}
	if got := fields["bar"].(string); got != "bar" {
		t.Fatalf("Fields: %s", got)
	}

	o.WithFieldPairs("foo", "foo", "bar").Trace("foo")
	if len(fields) != 2 {
		t.Fatalf("Fields: %v", fields)
	}
	if got := fields["foo"].(string); got != "foo" {
		t.Fatalf("Fields: %s", got)
	}
	if got := fields["bar"].(string); got != "" {
		t.Fatalf("Fields: %s", got)
	}

	o.WithFieldPairs(1, "foo").Trace("foo")
	if len(fields) != 1 {
		t.Fatalf("Fields: %v", fields)
	}
	if got := fields["1"].(string); got != "foo" {
		t.Fatalf("Fields: %s", got)
	}
}

func TestLogger_EnableHook(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)

	if o.EnableHook(false) == nil {
		t.Fatal("Logger.EnableHook(): return nil.")
	}

	var msg string
	o.AddHookFunc(GetAllLevels(), func(su Summary) error {
		msg = su.Message()
		return nil
	})

	o.Info("test")

	if msg != "" {
		t.Fatal(msg)
	}
}
