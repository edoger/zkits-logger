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
	"testing"
)

func TestDefaultTextFormatter(t *testing.T) {
	if DefaultTextFormatter() == nil {
		t.Fatal("DefaultTextFormatter(): nil")
	}
}

func TestDefaultQuoteTextFormatter(t *testing.T) {
	if DefaultQuoteTextFormatter() == nil {
		t.Fatal("DefaultQuoteTextFormatter(): nil")
	}
}

func TestNewTextFormatter(t *testing.T) {
	format := "[{name}][{time@2006-01-02 15:04:05}][{level@sc}] {caller} {message} {fields}"
	if f, err := NewTextFormatter(format, true); err != nil {
		t.Fatalf("NewTextFormatter(): error %s", err)
	} else {
		if f == nil {
			t.Fatal("NewTextFormatter(): nil")
		}
	}

	format = "hello"
	if _, err := NewTextFormatter(format, true); err == nil {
		t.Fatalf("NewTextFormatter(): no error with format: %q", format)
	}

	format = "hello {world}"
	if _, err := NewTextFormatter(format, true); err == nil {
		t.Fatalf("NewTextFormatter(): no error with format: %q", format)
	}
}

func TestMustNewTextFormatter(t *testing.T) {
	format := "[{name}][{time@2006-01-02 15:04:05}][{level@sc}] {caller} {message} {fields}"
	if MustNewTextFormatter(format, true) == nil {
		t.Fatal("MustNewTextFormatter(): nil")
	}

	format = "hello"
	defer func() {
		if recover() == nil {
			t.Fatalf("MustNewTextFormatter(): no panic with format : %q", format)
		}
	}()

	MustNewTextFormatter(format, true)
}

func TestTextFormatter_Format(t *testing.T) {
	l := New("test")
	l.SetFormatter(MustNewTextFormatter("{name} - {time@test} [{level@sc}] {caller@?} {message} {fields@?}", true))
	buf := new(bytes.Buffer)
	l.SetOutput(buf)

	l.WithField("foo", 1).WithField("bar", []byte("bar")).Info("test\n test")

	got := buf.String()
	want := "test - test [INF]  test\\n test bar=bar, foo=1\n"
	if got != want {
		t.Fatalf("TextFormatter.Format(): want %q, got %q", want, got)
	}

	buf.Reset()
	l.SetFormatter(MustNewTextFormatter("{name} - {time@test} [{level@sc}] {caller@?} {message} {fields@?}", false))

	l.WithField("foo", 1).WithField("bar", []byte("bar")).Info("test\n test")

	got = buf.String()
	want = "test - test [INF]  test\n test bar=bar, foo=1\n"
	if got != want {
		t.Fatalf("TextFormatter.Format(): want %q, got %q", want, got)
	}

	buf.Reset()
	l.SetFormatter(MustNewTextFormatter("{name} - {time@test} [{level@s}] {caller@?} {message} {fields@?}", true))

	l.Info("test")

	got = buf.String()
	want = "test - test [inf]  test \n"
	if got != want {
		t.Fatalf("TextFormatter.Format(): want %q, got %q", want, got)
	}

	buf.Reset()
	l.SetFormatter(MustNewTextFormatter("{name} - {time@test} [{level@c}] {caller@?} {message} {fields@?}", true))

	l.Info("test")

	got = buf.String()
	want = "test - test [INFO]  test \n"
	if got != want {
		t.Fatalf("TextFormatter.Format(): want %q, got %q", want, got)
	}

	buf.Reset()
	l.SetFormatter(MustNewTextFormatter("{name} - {time} [{level}] {caller@?} {message} {fields@?}", true))
	l.SetDefaultTimeFormat("test")
	l.Info("test")

	got = buf.String()
	want = "test - test [info]  test \n"
	if got != want {
		t.Fatalf("TextFormatter.Format(): want %q, got %q", want, got)
	}
}
