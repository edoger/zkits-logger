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

package internal

import (
	"bytes"
	"os"
	"testing"
)

func TestEmptyExitFunc(t *testing.T) {
	EmptyExitFunc(0) // nothing
}

func TestDefaultPanicFunc(t *testing.T) {
	defer func() {
		v := recover()
		if v == nil {
			t.Fatal("DefaultPanicFunc(): no panic")
		}
		s, ok := v.(string)
		if !ok {
			t.Fatalf("DefaultPanicFunc(): got %T - %v", v, v)
		}
		if s != "test" {
			t.Fatalf("DefaultPanicFunc(): got %s", s)
		}
	}()
	DefaultPanicFunc("test")
}

func TestEmptyPanicFunc(t *testing.T) {
	defer func() {
		v := recover()
		if v != nil {
			t.Fatalf("EmptyPanicFunc(): got %T - %v", v, v)
		}
	}()
	EmptyPanicFunc("test")
}

func TestEchoError(t *testing.T) {
	defer func() { ErrorWriter = os.Stderr }()

	buf := new(bytes.Buffer)
	ErrorWriter = buf

	EchoError("test-%d", 1)

	if got := buf.String(); got != "test-1\n" {
		t.Fatalf("EchoError(): got %q", got)
	}
}
