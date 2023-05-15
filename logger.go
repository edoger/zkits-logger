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
	stdlog "log"
	"os"
	"sync/atomic"
	"time"

	"github.com/edoger/zkits-logger/internal"
)

// Logger interface defines a standard logger.
type Logger interface {
	Log

	// GetLevel returns the current logger level.
	GetLevel() Level

	// SetLevel sets the current logger level.
	// When the given log level is invalid, this method does nothing.
	SetLevel(Level) Logger

	// SetLevelString sets the current logger level by string.
	SetLevelString(s string) error

	// ForceSetLevelString sets the current logger level by string.
	// When the given log level string is invalid, this method does nothing.
	ForceSetLevelString(s string) Logger

	// SetOutput sets the current logger output writer.
	// If the given writer is nil, os.Stdout is used.
	SetOutput(io.Writer) Logger

	// SetLevelOutput sets the current logger level output writer.
	// The level output writer is used to write log data of a given level.
	// If the given writer is nil, the level writer will be disabled.
	SetLevelOutput(Level, io.Writer) Logger

	// SetLevelsOutput sets the current logger levels output writer.
	// The level output writer is used to write log data of a given level.
	// If the given writer is nil, the levels writer will be disabled.
	SetLevelsOutput([]Level, io.Writer) Logger

	// SetOutputInterceptor sets the output interceptor for the current logger.
	// If the given interceptor is nil, the log data is written to the output writer.
	SetOutputInterceptor(func(Summary, io.Writer) (int, error)) Logger

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

	// SetFormatOutput sets the log format output.
	// After setting the format output, the format and output of the logger will be controlled by this structure,
	// and the bound log output and log level output will no longer be used.
	// If format output needs to be disabled, set to nil and the logger will back to the original behavior.
	SetFormatOutput(FormatOutput) Logger

	// SetDefaultTimeFormat sets the default log time format for the current logger.
	// If the given time format is empty string, internal.DefaultTimeFormat is used.
	SetDefaultTimeFormat(string) Logger

	// EnableCaller enables caller reporting on all levels of logs.
	EnableCaller(...int) Logger

	// EnableLevelCaller enables caller reporting on logs of a given level.
	EnableLevelCaller(Level, ...int) Logger

	// EnableLevelsCaller enables caller reporting on logs of the given levels.
	EnableLevelsCaller([]Level, ...int) Logger

	// AddHook adds the given log hook to the current logger.
	AddHook(Hook) Logger

	// AddHookFunc adds the given log hook function to the current logger.
	AddHookFunc([]Level, func(Summary) error) Logger

	// EnableHook enables or disables the log hook.
	EnableHook(bool) Logger

	// AsLog converts current Logger to Log instances, which is unidirectional.
	AsLog() Log

	// AsStandardLogger converts the current logger to a standard library logger instance.
	AsStandardLogger() *stdlog.Logger

	// SetStackPrefixFilter sets the call stack prefix filter rules.
	SetStackPrefixFilter(...string) Logger
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
	return Level(atomic.LoadUint32(&o.core.level))
}

// SetLevel sets the current logger level.
// When the given log level is invalid, this method does nothing.
func (o *logger) SetLevel(level Level) Logger {
	if level.IsValid() {
		atomic.StoreUint32(&o.core.level, uint32(level))
	}
	return o
}

// SetLevelString sets the current logger level by string.
func (o *logger) SetLevelString(s string) error {
	level, err := ParseLevel(s)
	if err != nil {
		return err
	}
	o.SetLevel(level)
	return nil
}

// ForceSetLevelString sets the current logger level by string.
// When the given log level string is invalid, this method does nothing.
func (o *logger) ForceSetLevelString(s string) Logger {
	if level, err := ParseLevel(s); err == nil {
		o.SetLevel(level)
	}
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
	if w == nil {
		delete(o.core.levelWriter, level)
	} else {
		o.core.levelWriter[level] = w
	}
	return o
}

// SetLevelsOutput sets the current logger levels output writer.
// The level output writer is used to write log data of a given level.
// If the given writer is nil, the levels writer will be disabled.
func (o *logger) SetLevelsOutput(levels []Level, w io.Writer) Logger {
	for i, j := 0, len(levels); i < j; i++ {
		o.SetLevelOutput(levels[i], w)
	}
	return o
}

// SetOutputInterceptor sets the output interceptor for the current logger.
// If the given interceptor is nil, the log data is written to the output writer.
func (o *logger) SetOutputInterceptor(f func(Summary, io.Writer) (int, error)) Logger {
	o.core.interceptor = f
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
	if formatter == nil {
		o.core.formatter = DefaultJSONFormatter()
	} else {
		o.core.formatter = formatter
	}
	return o
}

// SetFormatOutput sets the log format output.
// After setting the format output, the format and output of the logger will be controlled by this structure,
// and the bound log output and log level output will no longer be used.
// If format output needs to be disabled, set to nil and the logger will back to the original behavior.
func (o *logger) SetFormatOutput(fo FormatOutput) Logger {
	o.core.formatOutput = fo
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

// EnableCaller enables caller reporting on all levels of logs.
func (o *logger) EnableCaller(skip ...int) Logger {
	var n int
	if len(skip) > 0 && skip[0] > 0 {
		n = skip[0]
	}
	o.core.caller = internal.NewCallerReporter(n)
	return o
}

// EnableLevelCaller enables caller reporting on logs of a given level.
func (o *logger) EnableLevelCaller(level Level, skip ...int) Logger {
	var n int
	if len(skip) > 0 && skip[0] > 0 {
		n = skip[0]
	}
	o.core.levelCaller[level] = internal.NewCallerReporter(n)
	return o
}

// EnableLevelsCaller enables caller reporting on logs of the given levels.
func (o *logger) EnableLevelsCaller(levels []Level, skip ...int) Logger {
	for i, j := 0, len(levels); i < j; i++ {
		o.EnableLevelCaller(levels[i], skip...)
	}
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

// EnableHook enables or disables the log hook.
func (o *logger) EnableHook(ok bool) Logger {
	o.core.enableHooks = ok
	return o
}

// AsLog converts current Logger to Log instances, which is unidirectional.
func (o *logger) AsLog() Log {
	return &o.log
}

// AsStandardLogger converts the current logger to a standard library logger instance.
func (o *logger) AsStandardLogger() *stdlog.Logger {
	return stdlog.New(NewLevelWriter(InfoLevel, o.AsLog()), "", 0)
}

// SetStackPrefixFilter sets the call stack prefix filter rules.
func (o *logger) SetStackPrefixFilter(prefixes ...string) Logger {
	o.core.stackPrefixes = internal.FormatKnownStackPrefixes(prefixes...)
	return o
}
