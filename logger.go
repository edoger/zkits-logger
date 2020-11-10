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
	"io"
	"os"
	"time"

	"github.com/edoger/zkits-logger/internal"
)

type Logger interface {
	Log

	// GetLevel returns the current logger level.
	GetLevel() Level

	// SetLevel sets the current logger level.
	SetLevel(Level) Logger

	// SetOutput sets the current logger output writer.
	SetOutput(io.Writer) Logger

	SetNowFunc(func() time.Time) Logger
	SetExitFunc(func(int)) Logger
	SetPanicFunc(func(string)) Logger
	SetFormatter(Formatter) Logger
	AddHook(Hook) Logger
	AddHookFunc([]Level, func(Summary) error) Logger
}

// New creates a new Logger instance.
// By default, the logger level is TraceLevel and logs will be output to os.Stdout.
func New(name string) Logger {
	return &logger{log{core: newCore(name)}}
}

type logger struct {
	log
}

// GetLevel returns the current logger level.
func (o *logger) GetLevel() Level {
	return o.core.level
}

// SetLevel sets the current logger level.
func (o *logger) SetLevel(level Level) Logger {
	o.core.level = level
	return o
}

// SetOutput sets the current logger output writer.
func (o *logger) SetOutput(w io.Writer) Logger {
	if w == nil {
		o.core.writer = os.Stdout
	} else {
		o.core.writer = w
	}
	return o
}

func (o *logger) SetNowFunc(f func() time.Time) Logger {
	if f == nil {
		o.core.nowFunc = internal.DefaultNowFunc
	} else {
		o.core.nowFunc = f
	}
	return o
}

func (o *logger) SetExitFunc(f func(int)) Logger {
	if f == nil {
		o.core.exitFunc = internal.EmptyExitFunc
	} else {
		o.core.exitFunc = f
	}
	return o
}

func (o *logger) SetPanicFunc(f func(string)) Logger {
	if f == nil {
		o.core.panicFunc = internal.EmptyPanicFunc
	} else {
		o.core.panicFunc = f
	}
	return o
}

func (o *logger) SetFormatter(formatter Formatter) Logger {
	o.core.formatter = formatter
	return o
}

func (o *logger) AddHook(hook Hook) Logger {
	o.core.hooks.Add(hook)
	return o
}

func (o *logger) AddHookFunc(levels []Level, hook func(Summary) error) Logger {
	return o.AddHook(NewHookFromFunc(levels, hook))
}
