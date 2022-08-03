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

func TestLevel_ColorfulString(t *testing.T) {
	items := []struct {
		Given Level
		Want  string
	}{
		{PanicLevel, "\u001B[91mpanic\u001B[0m"},
		{FatalLevel, "\u001B[31mfatal\u001B[0m"},
		{ErrorLevel, "\u001B[95merror\u001B[0m"},
		{WarnLevel, "\u001B[93mwarn\u001B[0m"},
		{InfoLevel, "\u001B[92minfo\u001B[0m"},
		{DebugLevel, "\u001B[96mdebug\u001B[0m"},
		{TraceLevel, "\u001B[36mtrace\u001B[0m"},
		{Level(0), "unknown"},
	}

	for _, item := range items {
		if got := item.Given.ColorfulString(); got != item.Want {
			t.Fatalf("Level.ColorfulString(): want %v, got %v", item.Want, got)
		}
	}
}

func TestLevel_CapitalString(t *testing.T) {
	items := []struct {
		Given Level
		Want  string
	}{
		{PanicLevel, "PANIC"},
		{FatalLevel, "FATAL"},
		{ErrorLevel, "ERROR"},
		{WarnLevel, "WARN"},
		{InfoLevel, "INFO"},
		{DebugLevel, "DEBUG"},
		{TraceLevel, "TRACE"},
		{Level(0), "UNKNOWN"},
	}

	for _, item := range items {
		if got := item.Given.CapitalString(); got != item.Want {
			t.Fatalf("Level.CapitalString(): want %v, got %v", item.Want, got)
		}
	}
}

func TestLevel_ColorfulCapitalString(t *testing.T) {
	items := []struct {
		Given Level
		Want  string
	}{
		{PanicLevel, "\u001B[91mPANIC\u001B[0m"},
		{FatalLevel, "\u001B[31mFATAL\u001B[0m"},
		{ErrorLevel, "\u001B[95mERROR\u001B[0m"},
		{WarnLevel, "\u001B[93mWARN\u001B[0m"},
		{InfoLevel, "\u001B[92mINFO\u001B[0m"},
		{DebugLevel, "\u001B[96mDEBUG\u001B[0m"},
		{TraceLevel, "\u001B[36mTRACE\u001B[0m"},
		{Level(0), "UNKNOWN"},
	}

	for _, item := range items {
		if got := item.Given.ColorfulCapitalString(); got != item.Want {
			t.Fatalf("Level.ColorfulCapitalString(): want %v, got %v", item.Want, got)
		}
	}
}

func TestLevel_ShortString(t *testing.T) {
	items := []struct {
		Given Level
		Want  string
	}{
		{PanicLevel, "pnc"},
		{FatalLevel, "fat"},
		{ErrorLevel, "err"},
		{WarnLevel, "wan"},
		{InfoLevel, "inf"},
		{DebugLevel, "dbg"},
		{TraceLevel, "tac"},
		{Level(0), "uno"},
	}

	for _, item := range items {
		if got := item.Given.ShortString(); got != item.Want {
			t.Fatalf("Level.ShortString(): want %v, got %v", item.Want, got)
		}
	}
}

func TestLevel_ColorfulShortString(t *testing.T) {
	items := []struct {
		Given Level
		Want  string
	}{
		{PanicLevel, "\u001B[91mpnc\u001B[0m"},
		{FatalLevel, "\u001B[31mfat\u001B[0m"},
		{ErrorLevel, "\u001B[95merr\u001B[0m"},
		{WarnLevel, "\u001B[93mwan\u001B[0m"},
		{InfoLevel, "\u001B[92minf\u001B[0m"},
		{DebugLevel, "\u001B[96mdbg\u001B[0m"},
		{TraceLevel, "\u001B[36mtac\u001B[0m"},
		{Level(0), "uno"},
	}

	for _, item := range items {
		if got := item.Given.ColorfulShortString(); got != item.Want {
			t.Fatalf("Level.ColorfulShortString(): want %v, got %v", item.Want, got)
		}
	}
}

func TestLevel_ShortCapitalString(t *testing.T) {
	items := []struct {
		Given Level
		Want  string
	}{
		{PanicLevel, "PNC"},
		{FatalLevel, "FAT"},
		{ErrorLevel, "ERR"},
		{WarnLevel, "WAN"},
		{InfoLevel, "INF"},
		{DebugLevel, "DBG"},
		{TraceLevel, "TAC"},
		{Level(0), "UNO"},
	}

	for _, item := range items {
		if got := item.Given.ShortCapitalString(); got != item.Want {
			t.Fatalf("Level.ShortCapitalString(): want %v, got %v", item.Want, got)
		}
	}
}

func TestLevel_ColorfulShortCapitalString(t *testing.T) {
	items := []struct {
		Given Level
		Want  string
	}{
		{PanicLevel, "\u001B[91mPNC\u001B[0m"},
		{FatalLevel, "\u001B[31mFAT\u001B[0m"},
		{ErrorLevel, "\u001B[95mERR\u001B[0m"},
		{WarnLevel, "\u001B[93mWAN\u001B[0m"},
		{InfoLevel, "\u001B[92mINF\u001B[0m"},
		{DebugLevel, "\u001B[96mDBG\u001B[0m"},
		{TraceLevel, "\u001B[36mTAC\u001B[0m"},
		{Level(0), "UNO"},
	}

	for _, item := range items {
		if got := item.Given.ColorfulShortCapitalString(); got != item.Want {
			t.Fatalf("Level.ColorfulShortCapitalString(): want %v, got %v", item.Want, got)
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
		{"pnc", PanicLevel, false},
		{"fatal", FatalLevel, false},
		{"fat", FatalLevel, false},
		{"error", ErrorLevel, false},
		{"err", ErrorLevel, false},
		{"warn", WarnLevel, false},
		{"wan", WarnLevel, false},
		{"warning", WarnLevel, false},
		{"info", InfoLevel, false},
		{"inf", InfoLevel, false},
		{"echo", InfoLevel, false},
		{"debug", DebugLevel, false},
		{"dbg", DebugLevel, false},
		{"trace", TraceLevel, false},
		{"tac", TraceLevel, false},
		{"print", TraceLevel, false},
		{"unknown", 0, true},
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

func TestGetAllLevels(t *testing.T) {
	got := GetAllLevels()

	if len(got) != 7 {
		t.Fatalf("GetAllLevels(): %v", got)
	}
}
