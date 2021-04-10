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

func TestTextFormatterCaller(t *testing.T) {
	l := New("test")
	l.SetFormatter(DefaultTextFormatter())
	buf := new(bytes.Buffer)
	l.SetOutput(buf)
	l.SetDefaultTimeFormat("time")
	l.EnableCaller()

	l.Info("test-caller") // Line 30

	got := buf.String()
	want := "test:[time][info] test-caller text_formatter_caller_test.go:30\n"
	if got != want {
		t.Fatalf("TextFormatter.Format(): want %q, got %q", want, got)
	}
}
