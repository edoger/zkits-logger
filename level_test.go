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
	items := []struct {
		Given Level
		Want  string
	}{
		{PanicLevel, "panic"},
		{FatalLevel, "fatal"},
		{ErrorLevel, "error"},
		{WarnLevel, "warn"},
		{InfoLevel, "info"},
		{DebugLevel, "debug"},
		{TraceLevel, "trace"},
		{Level(0), "unknown"},
	}

	for _, item := range items {
		if got := item.Given.String(); got != item.Want {
			t.Fatalf("Level.String(): want %v, got %v", item.Want, got)
		}
	}
}

func TestLevel_IsValid(t *testing.T) {
	items := []struct {
		Given Level
		Want  bool
	}{
		{PanicLevel, true},
		{FatalLevel, true},
		{ErrorLevel, true},
		{WarnLevel, true},
		{InfoLevel, true},
		{DebugLevel, true},
		{TraceLevel, true},
		{Level(0), false},
	}

	for _, item := range items {
		if got := item.Given.IsValid(); got != item.Want {
			t.Fatalf("Level.IsValid(): want %v, got %v", item.Want, got)
		}
	}
}

func TestLevel_IsEnabled(t *testing.T) {
	level := InfoLevel
	items := []struct {
		Given Level
		Want  bool
	}{
		{PanicLevel, true},
		{FatalLevel, true},
		{ErrorLevel, true},
		{WarnLevel, true},
		{InfoLevel, true},
		{DebugLevel, false},
		{TraceLevel, false},
		{Level(0), false},
	}

	for _, item := range items {
		if got := level.IsEnabled(item.Given); got != item.Want {
			t.Fatalf("Level.IsEnabled(): want %v, got %v", item.Want, got)
		}
	}
}

func TestParseLevel(t *testing.T) {
	items := []struct {
		Given string
		Want  Level
		Erred bool
	}{
		{"panic", PanicLevel, false},
		{"fatal", FatalLevel, false},
		{"error", ErrorLevel, false},
		{"warn", WarnLevel, false},
		{"warning", WarnLevel, false},
		{"info", InfoLevel, false},
		{"debug", DebugLevel, false},
		{"trace", TraceLevel, false},
		{"unknown", Level(0), true},
	}

	for _, item := range items {
		got, err := ParseLevel(item.Given)
		if got != item.Want {
			t.Fatalf("ParseLevel(): want %v, got %v", item.Want, got)
		}
		if item.Erred {
			if err == nil {
				t.Fatal("ParseLevel(): no error")
			}
		} else {
			if err != nil {
				t.Fatalf("ParseLevel(): %s", err)
			}
		}
	}
}

func TestMustParseLevel(t *testing.T) {
	items := []struct {
		Given string
		Want  Level
	}{
		{"panic", PanicLevel},
		{"fatal", FatalLevel},
		{"error", ErrorLevel},
		{"warn", WarnLevel},
		{"warning", WarnLevel},
		{"info", InfoLevel},
		{"debug", DebugLevel},
		{"trace", TraceLevel},
	}

	for _, item := range items {
		if got := MustParseLevel(item.Given); got != item.Want {
			t.Fatalf("MustParseLevel(): want %v, got %v", item.Want, got)
		}
	}
}

func TestMustParseLevelPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("MustParseLevel() not panic")
		}
	}()

	MustParseLevel("unknown")
}
