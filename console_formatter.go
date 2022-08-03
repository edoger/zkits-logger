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

	"github.com/edoger/zkits-logger/internal"
)

// NewConsoleFormatter creates and returns an instance of the log console formatter.
// The console formatter is very similar to the text formatter. The only difference is that
// we output different console colors for different log levels, which is very useful when
// outputting logs from the console.
func NewConsoleFormatter() Formatter {
	return new(consoleFormatter)
}

// The built-in console formatter.
type consoleFormatter struct{}

// Format formats the given log entity into character data and writes it to the given buffer.
func (f *consoleFormatter) Format(e Entity, b *bytes.Buffer) (err error) {
	if name := e.Name(); name != "" {
		b.WriteString(name + " ")
	}
	if tm := e.TimeString(); tm != "" {
		b.WriteString("[" + tm + "]")
	}
	b.WriteString("[" + e.Level().ColorfulShortCapitalString() + "] ")
	b.WriteString(e.Message())
	if caller := e.Caller(); caller != "" {
		b.WriteString(" " + caller)
	}
	if fields := e.Fields(); len(fields) > 0 {
		b.WriteString(" " + internal.FormatFieldsToText(e.Fields()))
	}
	if stack := e.Stack(); len(stack) > 0 {
		// In the console, in order to be able to display the stack information better,
		// we will use a separate line to display the stack information.
		for i, j := 0, len(stack); i < j; i++ {
			b.WriteString("\n\t" + stack[i])
		}
	}
	b.WriteByte('\n')
	return
}
