// Copyright 2023 The ZKits Project Authors.
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
	"sort"
	"strings"
	"testing"
)

func TestGetStack(t *testing.T) {
	fn1 := func() []string {
		return GetStack(KnownStackPrefixes...)
	}
	fn2 := func() []string {
		return fn1()
	}
	fn3 := func() []string {
		return fn2()
	}

	got := fn3()
	gotString := strings.Join(got, "\n")
	fmt.Println(gotString)

	want := []string{
		"stack_test.go:26",
		"stack_test.go:29",
		"stack_test.go:32",
		"stack_test.go:35",
	}

	for _, s := range want {
		if !strings.Contains(gotString, s) {
			t.Fatalf("GetStack() want: %s", s)
		}
	}
}

func TestFormatKnownStackPrefixes(t *testing.T) {
	cloneKnownStackPrefixes := func() []string {
		r := make([]string, len(KnownStackPrefixes))
		copy(r, KnownStackPrefixes)
		return r
	}
	items := []struct {
		Given []string
		Want  []string
	}{
		{
			[]string{},
			KnownStackPrefixes,
		},
		{
			[]string{"foo"},
			append(cloneKnownStackPrefixes(), "foo"),
		},
		{
			[]string{"foo", "foo"},
			append(cloneKnownStackPrefixes(), "foo"),
		},
	}
	for i, item := range items {
		sort.Strings(item.Want)
		want := strings.Join(item.Want, ",")
		gotStrings := FormatKnownStackPrefixes(item.Given...)
		got := strings.Join(gotStrings, ",")
		if want != got {
			t.Fatalf("FormatKnownStackPrefixes() [%d]: %s", i+1, got)
		}
	}
}
