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
	"path/filepath"
	"runtime"
	"strconv"
)

// KnownCallerDepth is the internally known call stack depth.
const KnownCallerDepth = 5

// NewCallerReporter returns a CallerReporter instance.
func NewCallerReporter(skip int) *CallerReporter {
	if skip >= 0 && skip < 15 {
		return DefaultCallerReporter[skip]
	}
	return &CallerReporter{skip: skip}
}

// CallerReporter defines the log caller reporter.
type CallerReporter struct {
	skip int
}

// Equal determines whether the given skip is equal to the current caller reporter.
func (o *CallerReporter) Equal(skip int) bool {
	return o.skip == skip
}

// Skip gets the current skipped call stack depth.
func (o *CallerReporter) Skip() int {
	return o.skip
}

// GetCaller reports file and line number information about function invocations on
// the calling goroutine's stack.
func GetCaller(skipped int, long bool) string {
	if _, file, line, ok := runtime.Caller(skipped + KnownCallerDepth); ok {
		if base := filepath.Base(file); long {
			// Only the parent directory is added.
			return filepath.Join(filepath.Base(filepath.Dir(file)), base) + ":" + strconv.Itoa(line)
		} else {
			return base + ":" + strconv.Itoa(line)
		}
	}
	return "???:0"
}
