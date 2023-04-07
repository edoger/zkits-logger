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
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestMakeFields(t *testing.T) {
	if got := MakeFields(nil); got == nil || len(got) != 0 {
		t.Fatalf("MakeFields(): %v", got)
	}
	if got := MakeFields(map[string]interface{}{"a": "b"}); len(got) != 1 {
		t.Fatalf("MakeFields(): %v", got)
	} else {
		if fmt.Sprint(got["a"]) != "b" {
			t.Fatalf("MakeFields(): %v", got)
		}
	}
}

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

func TestStandardiseFieldsForJSONEncoder(t *testing.T) {
	src := map[string]interface{}{
		"foo": 1,
		"bar": errors.New("bar"),
	}
	dst := StandardiseFieldsForJSONEncoder(src)

	bs, err := json.Marshal(dst)
	if err != nil {
		t.Fatal(err)
	}

	want := `{"bar":"bar","foo":1}`
	got := string(bs)
	if want != got {
		t.Fatalf("StandardiseFieldsForJSONEncoder(): want %q, got %q", want, got)
	}
}

type testStringer struct {
	v string
}

func (o *testStringer) String() string {
	return o.v
}

func TestFormatFieldsToText(t *testing.T) {
	want := `bar=bar, baz=test, foo=1`
	got := FormatFieldsToText(map[string]interface{}{
		"foo": 1,
		"bar": []byte("bar"),
		"baz": &testStringer{v: "test"},
	})
	if want != got {
		t.Fatalf("FormatFieldsToText(): want %q, got %q", want, got)
	}
}

func TestFormatPairsToFields(t *testing.T) {
	got := FormatPairsToFields([]interface{}{
		"foo", "test",
		1, "test",
		&testStringer{v: "test"}, "test",
		&testStringer{v: "bar"},
	})
	want := map[string]interface{}{
		"foo":  "test",
		"1":    "test",
		"test": "test",
		"bar":  "",
	}
	if len(want) != len(got) {
		t.Fatalf("FormatPairsToFields(): %v", got)
	}
	for k, v := range want {
		// All value is string.
		if got[k].(string) != v.(string) {
			t.Fatalf("FormatPairsToFields(): %v", got)
		}
	}
}
