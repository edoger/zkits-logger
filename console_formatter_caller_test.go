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
	"fmt"
	"strconv"
	"testing"
)

func TestConsoleFormatter_Format_WithoutCaller(t *testing.T) {
	l := New("")
	l.SetFormatter(NewConsoleFormatter())
	l.SetDefaultTimeFormat("TIME")
	buf := new(bytes.Buffer)
	l.SetOutput(buf)

	l.WithField("foo", 1).WithCaller().Trace("test") // LINE 31
	fmt.Println(buf.String())
	want := "[TIME][\u001B[36mTAC\u001B[0m] test console_formatter_caller_test.go:31 foo=1\n"
	if want != buf.String() {
		t.Fatalf("NewConsoleFormatter().Format(): %s", strconv.Quote(buf.String()))
	}
}
