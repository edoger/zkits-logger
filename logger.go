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

// Logger interface defines a standard logger.
type Logger interface {
	Log

	// GetLevel returns the current logger level.
	GetLevel() Level

	// SetLevel sets the current logger level.
	SetLevel(Level) Logger

	// SetOutput sets the current logger output writer.
	// If the given writer is nil, os.Stdout is used.
	SetOutput(io.Writer) Logger

	// SetLevelOutput sets the current logger level output writer.
	// The level output writer is used to write log data of a given level.
	// If the given writer is nil, the level writer will be disabled.
	SetLevelOutput(level Level, w io.Writer) Logger

	// SetNowFunc sets the function that gets the current time.
	// If the given function is nil, time.Now is used.
	SetNowFunc(func() time.Time) Logger

	// SetExitFunc sets the exit function of the current logger.
	// If the given function is nil, the exit function is disabled.
	// The exit function is called automatically after the FatalLevel level log is recorded.
	// By default, the exit function we use is os.Exit.
	SetExitFunc(func(int)) Logger

	// SetPanicFunc sets the panic function of the current logger.
	// If the given function is nil, the panic function is disabled.
	// The panic function is called automatically after the PanicLevel level log is recorded.
	// By default, the panic function we use is func(s string) { panic(s) }.
	SetPanicFunc(func(string)) Logger

	// SetFormatter sets the log formatter for the current logger.
	// If the given log formatter is nil, we will record the log in JSON format.
	SetFormatter(Formatter) Logger

	// SetDefaultTimeFormat sets the default log time format for the current logger.
	// If the given time format is empty string, internal.DefaultTimeFormat is used.
	SetDefaultTimeFormat(string) Logger

	EnableCaller(...int) Logger
	EnableLevelCaller(Level, ...int) Logger

	// AddHook adds the given log hook to the current logger.
	AddHook(Hook) Logger

	// AddHookFunc adds the given log hook function to the current logger.
	AddHookFunc([]Level, func(Summary) error) Logger
}

// New creates a new Logger instance.
// By default, the logger level is TraceLevel and logs will be output to os.Stdout.
func New(name string) Logger {
	return &logger{log{core: newCore(name)}}
}

// The logger type is an implementation of the built-in logger interface.
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
// If the given writer is nil, os.Stdout is used.
func (o *logger) SetOutput(w io.Writer) Logger {
	if w == nil {
		o.core.writer = os.Stdout
	} else {
		o.core.writer = w
	}
	return o
}

// SetLevelOutput sets the current logger level output writer.
// The level output writer is used to write log data of a given level.
// If the given writer is nil, the level writer will be disabled.
func (o *logger) SetLevelOutput(level Level, w io.Writer) Logger {
	o.core.levelWriter[level] = w
	return o
}

// SetNowFunc sets the function that gets the current time.
// If the given function is nil, time.Now is used.
func (o *logger) SetNowFunc(f func() time.Time) Logger {
	if f == nil {
		o.core.nowFunc = internal.DefaultNowFunc
	} else {
		o.core.nowFunc = f
	}
	return o
}

// SetExitFunc sets the exit function of the current logger.
// If the given function is nil, the exit function is disabled.
// The exit function is called automatically after the FatalLevel level log is recorded.
// By default, the exit function we use is os.Exit.
func (o *logger) SetExitFunc(f func(int)) Logger {
	if f == nil {
		o.core.exitFunc = internal.EmptyExitFunc
	} else {
		o.core.exitFunc = f
	}
	return o
}

// SetPanicFunc sets the panic function of the current logger.
// If the given function is nil, the panic function is disabled.
// The panic function is called automatically after the PanicLevel level log is recorded.
// By default, the panic function we use is func(s string) { panic(s) }.
func (o *logger) SetPanicFunc(f func(string)) Logger {
	if f == nil {
		o.core.panicFunc = internal.EmptyPanicFunc
	} else {
		o.core.panicFunc = f
	}
	return o
}

// SetFormatter sets the log formatter for the current logger.
// If the given log formatter is nil, we will record the log in JSON format.
func (o *logger) SetFormatter(formatter Formatter) Logger {
	o.core.formatter = formatter
	return o
}

// SetDefaultTimeFormat sets the default log time format for the current logger.
// If the given time format is empty string, internal.DefaultTimeFormat is used.
func (o *logger) SetDefaultTimeFormat(format string) Logger {
	if format == "" {
		o.core.timeFormat = internal.DefaultTimeFormat
	} else {
		o.core.timeFormat = format
	}
	return o
}

func (o *logger) EnableCaller(skip ...int) Logger {
	var n int
	if len(skip) > 0 && skip[0] > 0 {
		n = skip[0]
	}
	o.core.caller = internal.NewCaller(n)
	return o
}

func (o *logger) EnableLevelCaller(level Level, skip ...int) Logger {
	var n int
	if len(skip) > 0 && skip[0] > 0 {
		n = skip[0]
	}
	o.core.levelCaller[level] = internal.NewCaller(n)
	return o
}

// AddHook adds the given log hook to the current logger.
func (o *logger) AddHook(hook Hook) Logger {
	o.core.hooks.Add(hook)
	return o
}

// AddHookFunc adds the given log hook function to the current logger.
func (o *logger) AddHookFunc(levels []Level, hook func(Summary) error) Logger {
	return o.AddHook(NewHookFromFunc(levels, hook))
}
