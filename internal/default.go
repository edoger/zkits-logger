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

package internal

import (
	"fmt"
	"io"
	"os"
	"time"
)

var (
	// DefaultExitFunc is the default exit function for all logger instances.
	DefaultExitFunc = os.Exit

	// EmptyExitFunc is the empty exit function for all logger instances.
	EmptyExitFunc = func(int) {}

	// DefaultPanicFunc is the default panic function for all logger instances.
	DefaultPanicFunc = func(s string) { panic(s) }

	// EmptyPanicFunc is the empty panic function for all logger instances.
	EmptyPanicFunc = func(string) {}

	// DefaultNowFunc is the default now function for all logger instances.
	DefaultNowFunc = time.Now

	// DefaultTimeFormat is the default time format for all logger instances.
	DefaultTimeFormat = time.RFC3339

	// DefaultCallerReporter is the default caller for all logger instances.
	// Creating these 15 instances in advance can already cover most cases.
	DefaultCallerReporter = [15]*CallerReporter{
		{skip: 0}, {skip: 1}, {skip: 2}, {skip: 3}, {skip: 4},
		{skip: 5}, {skip: 6}, {skip: 7}, {skip: 8}, {skip: 9},
		{skip: 10}, {skip: 11}, {skip: 12}, {skip: 13}, {skip: 14},
	}

	// ErrorWriter is used to output logger internal error messages.
	ErrorWriter io.Writer = os.Stderr
)

// EchoError writes the given logger internal error message to ErrorWriter.
func EchoError(format string, args ...interface{}) {
	_, _ = fmt.Fprintln(ErrorWriter, fmt.Sprintf(format, args...))
}
