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

	o.WithStack().WithStack().Info("stack")

	// goroutine 4 [running]:
	// github.com/edoger/zkits-logger.TestLog_WithStack(0xc0001421a0) At ... /logger_stack_test.go:33 +0x317
	// testing.tRunner(0xc0001421a0, 0x116cd98) At ... /src/testing/testing.go:1439 +0x102
	// created by testing.(*T).Run At ... /src/testing/testing.go:1486 +0x35f
	if s := w.String(); strings.TrimSpace(s) == "" {
		t.Fatal("Log.WithStack(): no stack data")
	}
}
