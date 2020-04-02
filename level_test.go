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
	"testing"
)

func TestLevel_String(t *testing.T) {
	for level, s := range iLevels {
		if got := level.String(); got != s {
			t.Fatalf("%s != %s", s, got)
		}
	}

	if got := Level(0).String(); got != "unknown" {
		t.Fatalf("unknown != %s", got)
	}
}

func TestLevel_IsValid(t *testing.T) {
	for level := range iLevels {
		if !level.IsValid() {
			t.Fatalf("Level(%d).IsValid() return false", level)
		}
	}

	if Level(0).IsValid() {
		t.Fatal("Level(0).IsValid() return true")
	}
}

func TestParseLevel(t *testing.T) {
	for level, s := range iLevels {
		got, err := ParseLevel(s)
		if err != nil {
			t.Fatal(err)
		}
		if got != level {
			t.Fatalf("%d != %d", got, level)
		}
	}

	if got, err := ParseLevel("unknown"); err == nil {
		t.Fatal(`ParseLevel("unknown") return nil error`)
	} else {
		if got != 0 {
			t.Fatalf(`ParseLevel("unknown") return %d`, got)
		}
	}
}

func TestMustParseLevel(t *testing.T) {
	for level, s := range iLevels {
		if got := MustParseLevel(s); got != level {
			t.Fatalf("%d != %d", got, level)
		}
	}

	call := func(f func()) (v interface{}) {
		defer func() { v = recover() }()
		f()
		return
	}

	v := call(func() { MustParseLevel("unknown") })
	if v == nil {
		t.Fatal(`MustParseLevel("unknown") not panic`)
	}
}
