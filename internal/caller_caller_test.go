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
	"testing"
)

func TestCallerReporter_GetCaller(t *testing.T) {
	r := NewCallerReporter(1)
	f1 := func() string { return r.GetCaller() }
	f2 := func() string { return f1() }
	f3 := func() string { return f2() }
	f4 := func() string { return f3() }
	f5 := func() string { return f4() }

	got := f5() // Line 29
	if want := "caller_caller_test.go:29"; got != want {
		t.Fatalf("CallerReporter.GetCaller(): got %q, want %q", got, want)
	}

	got = r.GetCaller()
	if want := "???:0"; got != want {
		t.Fatalf("CallerReporter.GetCaller(): got %q, want %q", got, want)
	}
}
