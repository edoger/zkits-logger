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
	"testing"
)

func TestFields_Clone(t *testing.T) {
	src := Fields{}

	if got := src.Clone(0); len(got) != 0 {
		t.Fatalf("Fields.Clone(): %v", got)
	}

	src["key"] = "foo"
	if got := src.Clone(0); len(got) != 1 {
		t.Fatalf("Fields.Clone(): %v", got)
	} else {
		if s := got["key"].(string); s != "foo" {
			t.Fatalf("Fields.Clone(): %v", s)
		}
	}
}

func TestFields_With(t *testing.T) {
	src := Fields{"key": "foo"}

	if got := src.With(nil); len(got) != 1 {
		t.Fatalf("Fields.With(): %v", got)
	} else {
		if s := got["key"].(string); s != "foo" {
			t.Fatalf("Fields.With(): %v", s)
		}
	}

	if got := src.With(map[string]interface{}{"key2": "bar"}); len(got) != 2 {
		t.Fatalf("Fields.With(): %v", got)
	} else {
		if s := got["key"].(string); s != "foo" {
			t.Fatalf("Fields.With(): %v", s)
		}
		if s := got["key2"].(string); s != "bar" {
			t.Fatalf("Fields.With(): %v", s)
		}
	}
}
