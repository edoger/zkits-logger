// Copyright 2023 The ZKits Project Authors.
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
	"errors"
	"io"
	"testing"
)

func TestUnimplementedFormatter(t *testing.T) {
	var v any = new(UnimplementedFormatter)
	if f, ok := v.(Formatter); ok {
		err := f.Format(nil, nil)
		if err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal("UnimplementedFormatter: not is Formatter")
	}
}

func TestFormatOutputFunc(t *testing.T) {
	f := FormatOutputFunc(func(_ Entity, _ *bytes.Buffer) (io.Writer, error) {
		return nil, errors.New("test")
	})
	if _, err := f.Format(nil, nil); err == nil {
		t.Fatal("FormatOutputFunc.Format(): nil error")
	} else {
		if err.Error() != "test" {
			t.Fatalf("FormatOutputFunc.Format(): %s", err)
		}
	}
}

func TestFormatOutput(t *testing.T) {
	w := new(bytes.Buffer)
	f := NewFormatOutput(FormatterFunc(func(e Entity, b *bytes.Buffer) error {
		b.WriteString(e.Message())
		return nil
	}), w)
	if f == nil {
		t.Fatal("NewFormatOutput(): nil")
	}

	o := New("test")
	o.SetFormatOutput(f)
	o.SetLevel(TraceLevel)

	o.Info("test")

	if got := w.String(); got != "test" {
		t.Fatalf("FormatOutput: %s", got)
	}
}

func TestLevelPriorityFormatOutput(t *testing.T) {
	w1 := new(bytes.Buffer)
	w2 := new(bytes.Buffer)
	f1 := NewFormatOutput(FormatterFunc(func(e Entity, b *bytes.Buffer) error {
		b.WriteString(e.Message())
		return nil
	}), w1)
	f2 := NewFormatOutput(FormatterFunc(func(e Entity, b *bytes.Buffer) error {
		b.WriteString(e.Message())
		return nil
	}), w2)
	if f1 == nil || f2 == nil {
		t.Fatal("NewFormatOutput(): nil")
	}
	f := NewLevelPriorityFormatOutput(f1, f2)
	if f == nil {
		t.Fatal("NewLevelPriorityFormatOutput(): nil")
	}

	o := New("test")
	o.SetFormatOutput(f)
	o.SetLevel(TraceLevel)

	o.Info("info")
	o.Error("error")

	if got := w1.String(); got != "error" {
		t.Fatalf("LevelPriorityFormatOutput: %s", got)
	}
	if got := w2.String(); got != "info" {
		t.Fatalf("LevelPriorityFormatOutput: %s", got)
	}
}
