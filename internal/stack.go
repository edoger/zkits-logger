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

package internal

import (
	"bytes"
	"runtime"
	"strings"
)

// KnownStack defines the stack information when calling the GetStack method.
// We need to check this information from the call stack, because they are always fixed.
const KnownStack = "github.com/edoger/zkits-logger/internal.GetStack()"

// KnownStackPrefix defines the stack information prefix when calling the GetStack method.
// They look like this (excluding local path information):
//   - github.com/edoger/zkits-logger.(*log).record(...)
//   - github.com/edoger/zkits-logger.(*log).log(...)
//   - github.com/edoger/zkits-logger.(*log).METHOD(...)
// We need to check this information from the call stack, because they are always fixed.
const KnownStackPrefix = "github.com/edoger/zkits-logger.(*log)."

// GetStack returns the current coroutine call stack information.
// This method call is very expensive, we format the stack information returned by the system
// and exclude the internal call stack information (they are always unchanged).
func GetStack() (r []string) {
	// We use a 8KB buffer to read the call stack information, which is sufficient in most
	// cases, but it is not excluded that the size of the stack information exceeds it.
	// The reason why we do not use a larger buffer is because in almost In all scenarios,
	// the information at the bottom of the stack is enough to assist us in debugging and
	// capturing key information, and there is absolutely no need for us to obtain the
	// information at the top of the stack.
	buf := make([]byte, 1024*8)
	n := runtime.Stack(buf, false)
	// The called function name and file location are always paired, and we need to determine
	// whether we need to discard the file location information to ensure complete exclusion
	// of internal call stack information.
	skip := false
	for _, q := range bytes.Split(buf[:n], []byte{'\n'}) {
		if len(q) == 0 {
			continue
		}
		if q[0] != '\t' {
			// Determine whether this information is internal call stack information, and if
			// we find them, exclude the following file location information.
			if s := string(q); s == KnownStack || strings.HasPrefix(s, KnownStackPrefix) {
				skip = true
			} else {
				r = append(r, s)
			}
			continue
		}
		if skip {
			skip = false
			continue
		}
		if i := len(r); i == 0 {
			r = append(r, string(q[1:]))
		} else {
			// We combine the name of the function being called and the file location it is in
			// into a single message, which will be easier to read.
			r[i-1] += " At " + string(q[1:])
		}
	}
	return
}
