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

func (o testStringer) String() string {
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

func TestToString(t *testing.T) {
	s := "s"
	n := 1
	b := true
	f := 1.5
	f32 := float32(1.5)
	iu64 := uint64(1)

	ps := &s
	pn := &n
	pb := &b
	pf := &f
	pf32 := &f32
	piu64 := &iu64

	psr := &testStringer{"s"}
	psr2 := &psr

	bs := []byte("bs")

	nilError := error(nil)
	nilSr := fmt.Stringer(nil)

	type aliasInt int
	type aliasString string
	type aliasBool bool
	type aliasError error
	type aliasStringer fmt.Stringer

	ai := aliasInt(1)
	as := aliasString("s")
	ab := aliasBool(true)
	ae := aliasError(errors.New("err"))
	asr1 := aliasStringer(&testStringer{"s"})
	asr2 := aliasStringer(testStringer{"s"})

	nilAe := aliasError(nil)
	nilAsr := aliasStringer(nil)

	items := []struct {
		Given interface{}
		Want  string
	}{
		{"s", "s"},
		{1, "1"}, // int
		{int8(1), "1"},
		{int16(1), "1"},
		{int32(1), "1"},
		{int64(1), "1"},
		{uint(1), "1"},
		{uint8(1), "1"},
		{uint16(1), "1"},
		{uint32(1), "1"},
		{uint64(1), "1"},
		{uint8(1), "1"},
		{1.5, "1.5"}, // float64
		{float32(1.5), "1.5"},
		{[]byte("b"), "b"},
		{errors.New("err"), "err"},
		{error(nil), ""},
		{nilError, ""},
		{&nilError, ""},
		{fmt.Stringer(nil), ""},
		{nilSr, ""},
		{&nilSr, ""},
		{testStringer{"s"}, "s"},
		{&testStringer{"s"}, "s"},
		{&psr, "s"},
		{psr2, "s"},
		{true, "true"},
		{ps, "s"},      // *string
		{pn, "1"},      // *int
		{pb, "true"},   // *bool
		{pf, "1.5"},    // *float64
		{pf32, "1.5"},  // *float32
		{piu64, "1"},   // *uint64
		{&ps, "s"},     // **string
		{&pn, "1"},     // **int
		{&pb, "true"},  // **bool
		{&pf, "1.5"},   // **float64
		{&pf32, "1.5"}, // **float32
		{&piu64, "1"},  // **uint64
		{&bs, "bs"},    // *[]byte
		{[]string{"s", "s"}, fmt.Sprint([]string{"s", "s"})},
		{nil, ""},
		{aliasInt(1), "1"},
		{aliasString("s"), "s"},
		{aliasBool(true), "true"},
		{aliasError(nil), ""},
		{aliasError(errors.New("err")), "err"},
		{aliasStringer(nil), ""},
		{aliasStringer(testStringer{"s"}), "s"},
		{aliasStringer(&testStringer{"s"}), "s"},
		{&ai, "1"},
		{&as, "s"},
		{&ab, "true"},
		{&ae, "err"},
		{&asr1, "s"},
		{&asr2, "s"},
		{&nilAe, ""},
		{&nilAsr, ""},
	}

	for i, item := range items {
		if got := ToString(item.Given); got != item.Want {
			t.Fatalf("ToString(): [%d] want %q got %q", i, item.Want, got)
		}
	}
}

type testPanicError struct{}

func (testPanicError) Error() string {
	panic("testPanicError")
}

type testPanicStringer struct{}

func (testPanicStringer) String() string {
	panic("testPanicStringer")
}

func TestStandardiseFieldsForJSONEncoder_WithPanicError(t *testing.T) {
	src := map[string]interface{}{
		"err": new(testPanicError),
	}
	dst := StandardiseFieldsForJSONEncoder(src)

	bs, err := json.Marshal(dst)
	if err != nil {
		t.Fatal(err)
	}

	want := `{"err":"!!PANIC(error.Error)"}`
	got := string(bs)
	if want != got {
		t.Fatalf("StandardiseFieldsForJSONEncoder(): want %q, got %q", want, got)
	}
}

func TestFormatFieldsToText_WithPanicError(t *testing.T) {
	want := `err=!!PANIC(error.Error)`
	got := FormatFieldsToText(map[string]interface{}{
		"err": new(testPanicError),
	})
	if want != got {
		t.Fatalf("FormatFieldsToText(): want %q, got %q", want, got)
	}
}

func TestFormatFieldsToText_WithPanicStringer(t *testing.T) {
	want := `stringer=!!PANIC(fmt.Stringer.String)`
	got := FormatFieldsToText(map[string]interface{}{
		"stringer": new(testPanicStringer),
	})
	if want != got {
		t.Fatalf("FormatFieldsToText(): want %q, got %q", want, got)
	}
}

func TestFormatPairsToFields_WithPanicError(t *testing.T) {
	got := FormatPairsToFields([]interface{}{
		new(testPanicError), "test",
	})
	want := map[string]interface{}{
		"!!PANIC(error.Error)": "test",
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

func TestFormatPairsToFields_WithPanicStringer(t *testing.T) {
	got := FormatPairsToFields([]interface{}{
		new(testPanicStringer), "test",
	})
	want := map[string]interface{}{
		"!!PANIC(fmt.Stringer.String)": "test",
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
