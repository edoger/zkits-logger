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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func TestLogEntityAndSummary(t *testing.T) {
	now := time.Now()
	var o Summary = &logEntity{
		name:       "test",
		time:       now,
		timeFormat: time.RFC3339,
		level:      InfoLevel,
		message:    "foo",
		fields:     map[string]interface{}{"key": "foo"},
		caller:     "foo.go:1",
		stack:      []string{"stack"},
	}

	o.(*logEntity).buffer.WriteString("test")

	if got := o.Name(); got != "test" {
		t.Fatalf("Summary.Name(): %s", got)
	}
	if got := o.Time(); !got.Equal(now) {
		t.Fatalf("Summary.Time(): %s", got)
	}
	if got := o.TimeString(); got != now.Format(time.RFC3339) {
		t.Fatalf("Summary.TimeString(): %s", got)
	}
	// empty time string
	o.(*logEntity).timeFormat = ""
	if got := o.TimeString(); got != "" {
		t.Fatalf("Summary.TimeString(): %s", got)
	}

	if got := o.Level(); got != InfoLevel {
		t.Fatalf("Summary.Level(): %s", got.String())
	}
	if got := o.Message(); got != "foo" {
		t.Fatalf("Summary.Level(): %s", got)
	}
	if got := o.HasFields(); got != true {
		t.Fatalf("Summary.HasFields(): %v", got)
	}
	if got := o.Fields(); fmt.Sprint(got) != fmt.Sprint(map[string]interface{}{"key": "foo"}) {
		t.Fatalf("Summary.Fields(): %v", got)
	}
	if got := o.HasContext(); got == true {
		t.Fatalf("Summary.HasContext(): %v", got)
	}
	if got := o.Context(); got == nil {
		t.Fatal("Summary.Context(): nil")
	}
	if got := o.HasCaller(); got == false {
		t.Fatalf("Summary.HasCaller(): %v", got)
	}
	if got := o.Caller(); got != "foo.go:1" {
		t.Fatalf("Summary.Caller(): %s", got)
	}
	if got := o.HasStack(); got == false {
		t.Fatalf("Summary.HasStack(): %v", got)
	}
	if got := o.Stack(); len(got) != 1 || got[0] != "stack" {
		t.Fatalf("Summary.Stack(): %s", got)
	}
	if got := o.Buffer(); got == nil || got.String() != "test" {
		if got == nil {
			t.Fatal("Summary.Buffer(): nil")
		} else {
			t.Fatalf("Summary.Bytes(): %s", got.String())
		}
	}
	if got := o.Bytes(); !bytes.Equal(got, []byte("test")) {
		t.Fatalf("Summary.Bytes(): %s", string(got))
	}
	if got := o.String(); got != "test" {
		t.Fatalf("Summary.String(): %s", got)
	}
	if got := o.Size(); got != 4 {
		t.Fatalf("Summary.Size(): %d", got)
	}

	if o.Clone() == nil {
		t.Fatalf("Summary.Clone(): nil")
	}
	if o.CloneWithContext(context.Background()) == nil {
		t.Fatalf("Summary.CloneWithContext(): nil")
	}

	if got, err := ioutil.ReadAll(o); err == nil {
		if !bytes.Equal(got, []byte("test")) {
			t.Fatalf("Summary.Read(): %s", string(got))
		}
	} else {
		t.Fatalf("Summary.Read(): %s", err)
	}
	// empty buffer
	if o.Clone() == nil {
		t.Fatalf("Summary.Clone(): nil")
	}
	// empty buffer
	if o.CloneWithContext(context.Background()) == nil {
		t.Fatalf("Summary.CloneWithContext(): nil")
	}
}
