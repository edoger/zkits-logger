// Copyright 2022 The ZKits Project Authors.
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
	"strings"
	"testing"
)

func TestLog_WithStack(t *testing.T) {
	w := new(bytes.Buffer)
	o := New("test")
	o.SetOutput(w)
	o.SetLevel(TraceLevel)
	o.SetFormatter(FormatterFunc(func(e Entity, b *bytes.Buffer) error {
		b.WriteString(strings.Join(e.Stack(), "\n"))
		return nil
	}))

	l := o.WithStack().WithStack()
	l.Info("stack")

	// goroutine 18 [running]:
	// github.com/edoger/zkits-logger.TestLog_WithStack(0xc00010a820) At ... logger_stack_test.go:35 +0x323
	// testing.tRunner(0xc00010a820, 0x118a118) At ... testing.go:1576 +0x10b
	// created by testing.(*T).Run At ... testing.go:1629 +0x3ea
	if s := w.String(); strings.TrimSpace(s) == "" {
		t.Fatal("Log.WithStack(): no stack data")
	}
	fmt.Println(w.String())
}
