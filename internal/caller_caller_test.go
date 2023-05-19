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
	"strings"
	"testing"
)

func TestGetCaller(t *testing.T) {
	f1 := func() string { return GetCaller(1, false) }
	f2 := func() string { return f1() }
	f3 := func() string { return f2() }
	f4 := func() string { return f3() }
	f5 := func() string { return f4() }

	got := f5() // Line 29
	if want := "caller_caller_test.go:29"; got != want {
		t.Fatalf("GetCaller(): got %q, want %q", got, want)
	}

	got = GetCaller(0, false)
	if want := "???:0"; got != want {
		t.Fatalf("GetCaller(): got %q, want %q", got, want)
	}
}

func TestGetCaller_Long(t *testing.T) {
	f1 := func() string { return GetCaller(1, true) }
	f2 := func() string { return f1() }
	f3 := func() string { return f2() }
	f4 := func() string { return f3() }
	f5 := func() string { return f4() }

	got := f5() // Line 47
	if want := "caller_caller_test.go:47"; !strings.HasSuffix(got, want) {
		t.Fatalf("GetCaller(): got %q", got)
	}

	got = GetCaller(0, true)
	if want := "???:0"; got != want {
		t.Fatalf("GetCaller(): got %q, want %q", got, want)
	}
}
