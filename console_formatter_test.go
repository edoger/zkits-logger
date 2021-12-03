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
	"strings"
	"testing"
)

func TestNewConsoleFormatter(t *testing.T) {
	if f := NewConsoleFormatter(); f == nil {
		t.Fatal("NewConsoleFormatter(): return nil.")
	}
}

func TestConsoleFormatter_Format(t *testing.T) {
	l := New("CONSOLE")
	l.SetFormatter(NewConsoleFormatter())
	l.SetDefaultTimeFormat("TIME")
	l.SetExitFunc(nil)
	l.SetPanicFunc(nil)
	buf := new(bytes.Buffer)
	l.SetOutput(buf)

	l.WithField("foo", 1).Trace("test")
	l.WithField("foo", 1).Debug("test")
	l.WithField("foo", 1).Info("test")
	l.WithField("foo", 1).Warn("test")
	l.WithField("foo", 1).Error("test")
	l.WithField("foo", 1).Fatal("test")
	l.WithField("foo", 1).Panic("test")
	fmt.Println(buf.String())
	lines := []string{
		"CONSOLE [TIME][\u001B[36mTAC\u001B[0m] test foo=1",
		"CONSOLE [TIME][\u001B[96mDBG\u001B[0m] test foo=1",
		"CONSOLE [TIME][\u001B[92mINF\u001B[0m] test foo=1",
		"CONSOLE [TIME][\u001B[93mWAN\u001B[0m] test foo=1",
		"CONSOLE [TIME][\u001B[95mERR\u001B[0m] test foo=1",
		"CONSOLE [TIME][\u001B[31mFAT\u001B[0m] test foo=1",
		"CONSOLE [TIME][\u001B[91mPNC\u001B[0m] test foo=1",
	}
	if want := strings.Join(lines, "\n") + "\n"; want != buf.String() {
		t.Fatalf("NewConsoleFormatter().Format(): %s", strconv.Quote(buf.String()))
	}
}

func TestConsoleFormatter_Format_WithoutName(t *testing.T) {
	l := New("")
	l.SetFormatter(NewConsoleFormatter())
	l.SetDefaultTimeFormat("TIME")
	buf := new(bytes.Buffer)
	l.SetOutput(buf)

	l.WithField("foo", 1).Trace("test")
	fmt.Println(buf.String())
	want := "[TIME][\u001B[36mTAC\u001B[0m] test foo=1\n"
	if want != buf.String() {
		t.Fatalf("NewConsoleFormatter().Format(): %s", strconv.Quote(buf.String()))
	}
}
