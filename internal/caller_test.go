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

func TestNewCallerReporter(t *testing.T) {
	if NewCallerReporter(0) == nil {
		t.Fatal("NewCallerReporter(0): nil")
	}
	if NewCallerReporter(1) == nil {
		t.Fatal("NewCallerReporter(1): nil")
	}
}

func TestCallerReporter_Equal(t *testing.T) {
	r := NewCallerReporter(1)
	if r.Equal(0) {
		t.Fatal("CallerReporter.Equal(0): true")
	}
	if !r.Equal(1) {
		t.Fatal("CallerReporter.Equal(1): false")
	}
}
