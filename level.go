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
	"fmt"
	"strings"
)

const (
	// PanicLevel indicates a very serious event at the highest level.
	PanicLevel Level = iota + 1

	// FatalLevel indicates that an event occurred that the application
	// cannot continue to run.
	FatalLevel

	// ErrorLevel indicates that an event occurred within the application
	// but does not affect continued operation.
	ErrorLevel

	// WarnLevel indicates that a noteworthy event has occurred inside
	// the application.
	WarnLevel

	// InfoLevel represents some general information events.
	InfoLevel

	// DebugLevel represents some informational events for debugging, and
	// very verbose logging.
	DebugLevel

	// TraceLevel represents the finest granular information event.
	TraceLevel
)

// Level is the level of the log.
type Level uint32

// All supported log levels.
var allLevels = map[Level]string{
	PanicLevel: "panic", FatalLevel: "fatal",
	ErrorLevel: "error", WarnLevel: "warn", InfoLevel: "info",
	DebugLevel: "debug", TraceLevel: "trace",
}

// String returns the string form of the current level.
// If the log level is not supported, always returns "unknown".
func (level Level) String() string {
	if s, found := allLevels[level]; found {
		return s
	}
	return "unknown"
}

// IsValid determines whether the current level is valid.
func (level Level) IsValid() bool {
	return level <= TraceLevel && level >= PanicLevel
}

// IsEnabled returns whether the given level is included in the current level.
func (level Level) IsEnabled(l Level) bool {
	return l <= level && l > 0
}

// ParseLevel parses the log level from the given string.
func ParseLevel(s string) (Level, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "trace":
		return TraceLevel, nil
	}
	// A level zero value is not a supported level.
	return 0, fmt.Errorf(`invalid log level string "%s"`, s)
}

// MustParseLevel parses the log level from the given string.
// If the given string is invalid, it will panic.
func MustParseLevel(s string) Level {
	level, err := ParseLevel(s)
	if err != nil {
		panic(err)
	}
	return level
}
