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

func TestDefaultJSONFormatter(t *testing.T) {
	if DefaultJSONFormatter() == nil {
		t.Fatal("DefaultJSONFormatter(): nil")
	}
}

func TestNewJSONFormatter(t *testing.T) {
	f, err := NewJSONFormatter(map[string]string{"message": "msg"}, true)
	if err != nil {
		t.Fatalf("NewJSONFormatter(): error %s", err)
	}
	if f == nil {
		t.Fatal("NewJSONFormatter(): nil")
	}

	_, err = NewJSONFormatter(map[string]string{"hello": "hello"}, true)
	if err == nil {
		t.Fatal("NewJSONFormatter(): no error")
	}
}

func TestMustNewJSONFormatter(t *testing.T) {
	if MustNewJSONFormatter(map[string]string{"message": "msg"}, true) == nil {
		t.Fatal("MustNewJSONFormatter(): nil")
	}

	defer func() {
		if recover() == nil {
			t.Fatal("MustNewJSONFormatter(): no panic")
		}
	}()

	MustNewJSONFormatter(map[string]string{"hello": "hello"}, true)
}

func TestJSONFormatter_Format(t *testing.T) {
	l := New("test")
	l.SetFormatter(MustNewJSONFormatter(map[string]string{"message": "msg"}, true))
	buf := new(bytes.Buffer)
	l.SetOutput(buf)
	l.SetDefaultTimeFormat("test")

	l.Info("test")

	got := buf.String()
	want := `{"caller":"","fields":{},"level":"info","msg":"test","name":"test","time":"test"}` + "\n"
	if got != want {
		t.Fatalf("JSONFormatter.Format(): want %q, got %q", want, got)
	}
}
