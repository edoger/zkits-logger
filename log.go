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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/edoger/zkits-logger/internal"
)

// Log interface defines an extensible log.
type Log interface {
	// Name returns the logger name.
	Name() string

	// WithField adds the given extended data to the log.
	WithField(string, interface{}) Log

	// WithError adds the given error to the log.
	// This method is relative to WithField("error", error).
	WithError(error) Log

	// WithFields adds the given multiple extended data to the log.
	WithFields(map[string]interface{}) Log

	// WithContext adds the given context to the log.
	WithContext(context.Context) Log

	// WithCaller forces the caller report of the current log to be enabled.
	WithCaller(...int) Log

	// Log uses the given parameters to record a log of the specified level.
	// If the given log level is PanicLevel, the given panic function will be
	// called automatically after logging is completed.
	// If the given log level is FatalLevel, the given exit function will be
	// called automatically after logging is completed.
	// If the given log level is invalid, the log will be discarded.
	Log(Level, ...interface{})

	// Logln uses the given parameters to record a log of the specified level.
	// If the given log level is PanicLevel, the given panic function will be
	// called automatically after logging is completed.
	// If the given log level is FatalLevel, the given exit function will be
	// called automatically after logging is completed.
	// If the given log level is invalid, the log will be discarded.
	Logln(Level, ...interface{})

	// Logf uses the given parameters to record a log of the specified level.
	// If the given log level is PanicLevel, the given panic function will be
	// called automatically after logging is completed.
	// If the given log level is FatalLevel, the given exit function will be
	// called automatically after logging is completed.
	// If the given log level is invalid, the log will be discarded.
	Logf(Level, string, ...interface{})

	// Trace uses the given parameters to record a TraceLevel log.
	Trace(...interface{})

	// Traceln uses the given parameters to record a TraceLevel log.
	Traceln(...interface{})

	// Tracef uses the given parameters to record a TraceLevel log.
	Tracef(string, ...interface{})

	// Print uses the given parameters to record a TraceLevel log.
	Print(...interface{})

	// Println uses the given parameters to record a TraceLevel log.
	Println(...interface{})

	// Printf uses the given parameters to record a TraceLevel log.
	Printf(string, ...interface{})

	// Debug uses the given parameters to record a DebugLevel log.
	Debug(...interface{})

	// Debugln uses the given parameters to record a DebugLevel log.
	Debugln(...interface{})

	// Debugf uses the given parameters to record a DebugLevel log.
	Debugf(string, ...interface{})

	// Info uses the given parameters to record a InfoLevel log.
	Info(...interface{})

	// Infoln uses the given parameters to record a InfoLevel log.
	Infoln(...interface{})

	// Infof uses the given parameters to record a InfoLevel log.
	Infof(string, ...interface{})

	// Echo uses the given parameters to record a InfoLevel log.
	Echo(...interface{})

	// Echoln uses the given parameters to record a InfoLevel log.
	Echoln(...interface{})

	// Echof uses the given parameters to record a InfoLevel log.
	Echof(string, ...interface{})

	// Warn uses the given parameters to record a WarnLevel log.
	Warn(...interface{})

	// Warnln uses the given parameters to record a WarnLevel log.
	Warnln(...interface{})

	// Warnf uses the given parameters to record a WarnLevel log.
	Warnf(string, ...interface{})

	// Warning uses the given parameters to record a WarnLevel log.
	Warning(...interface{})

	// Warningln uses the given parameters to record a WarnLevel log.
	Warningln(...interface{})

	// Warningln uses the given parameters to record a WarnLevel log.
	Warningf(string, ...interface{})

	// Error uses the given parameters to record a ErrorLevel log.
	Error(...interface{})

	// Errorln uses the given parameters to record a ErrorLevel log.
	Errorln(...interface{})

	// Errorf uses the given parameters to record a ErrorLevel log.
	Errorf(string, ...interface{})

	// Fatal uses the given parameters to record a FatalLevel log.
	// After the log record is completed, the system will automatically call
	// the exit function given in advance.
	Fatal(...interface{})

	// Fatalln uses the given parameters to record a FatalLevel log.
	// After the log record is completed, the system will automatically call
	// the exit function given in advance.
	Fatalln(...interface{})

	// Fatalf uses the given parameters to record a FatalLevel log.
	// After the log record is completed, the system will automatically call
	// the exit function given in advance.
	Fatalf(string, ...interface{})

	// Panic uses the given parameters to record a PanicLevel log.
	// After the log record is completed, the system will automatically call
	// the panic function given in advance.
	Panic(...interface{})

	// Panicln uses the given parameters to record a PanicLevel log.
	// After the log record is completed, the system will automatically call
	// the panic function given in advance.
	Panicln(...interface{})

	// Panicf uses the given parameters to record a PanicLevel log.
	// After the log record is completed, the system will automatically call
	// the panic function given in advance.
	Panicf(string, ...interface{})
}

// The core type defines the collection of shared attributes within the log,
// and each independent Logger shares the same core instance.
type core struct {
	name        string
	level       Level
	formatter   Formatter
	writer      io.Writer
	levelWriter map[Level]io.Writer
	pool        sync.Pool
	hooks       HookBag
	timeFormat  string
	nowFunc     func() time.Time
	exitFunc    func(int)
	panicFunc   func(string)
	caller      *internal.CallerReporter
	levelCaller map[Level]*internal.CallerReporter
	interceptor func(Summary, io.Writer) (int, error)
}

// Create a new core instance and bind the logger name.
func newCore(name string) *core {
	return &core{
		name:        name,
		level:       TraceLevel,
		writer:      os.Stdout,
		levelWriter: make(map[Level]io.Writer),
		pool:        sync.Pool{New: func() interface{} { return new(logEntity) }},
		hooks:       NewHookBag(),
		timeFormat:  internal.DefaultTimeFormat,
		nowFunc:     internal.DefaultNowFunc,
		exitFunc:    internal.DefaultExitFunc,
		panicFunc:   internal.DefaultPanicFunc,
		levelCaller: make(map[Level]*internal.CallerReporter),
	}
}

// Get a log entity from the pool and initialize it.
func (c *core) getEntity(l *log, level Level, message, caller string) *logEntity {
	o := c.pool.Get().(*logEntity)

	o.name = c.name
	o.time = c.nowFunc()
	o.timeFormat = c.timeFormat
	o.level = level
	o.message = message
	o.fields = l.fields
	o.ctx = l.ctx
	o.caller = caller

	return o
}

// Clean up and recycle the given log entity.
func (c *core) putEntity(o *logEntity) {
	// If the log size exceeds 1KB, we need to discard this buffer to
	// free memory faster.
	if o.buffer.Cap() > 1024 {
		o.buffer = bytes.Buffer{}
	} else {
		o.buffer.Reset()
	}

	o.name = ""
	o.timeFormat = ""
	o.message = ""
	o.fields = nil
	o.ctx = nil
	o.caller = ""

	c.pool.Put(o)
}

// Internal implementation of the Log interface.
type log struct {
	core   *core
	ctx    context.Context
	fields internal.Fields
	caller *internal.CallerReporter
}

// Name returns the logger name.
func (o *log) Name() string {
	return o.core.name
}

// WithField adds the given extended data to the log.
func (o *log) WithField(key string, value interface{}) Log {
	r := &log{core: o.core, fields: o.fields.Clone(1), ctx: o.ctx, caller: o.caller}
	r.fields[key] = value
	return r
}

// WithError adds the given error to the log.
// This method is relative to WithField("error", error).
func (o *log) WithError(err error) Log {
	return o.WithField("error", err)
}

// WithFields adds the given multiple extended data to the log.
func (o *log) WithFields(fields map[string]interface{}) Log {
	return &log{core: o.core, fields: o.fields.With(fields), ctx: o.ctx, caller: o.caller}
}

// WithContext adds the given context to the log.
func (o *log) WithContext(ctx context.Context) Log {
	return &log{core: o.core, fields: o.fields.Clone(0), ctx: ctx, caller: o.caller}
}

// WithCaller forces the caller report of the current log to be enabled.
func (o *log) WithCaller(skip ...int) Log {
	var n int
	if len(skip) > 0 && skip[0] > 0 {
		n = skip[0]
	}
	// If the caller is equaled, we don't need to create a new log instance.
	if o.caller != nil && o.caller.Equal(n) {
		return o
	}
	return &log{core: o.core, fields: o.fields.Clone(0), ctx: o.ctx, caller: internal.NewCallerReporter(n)}
}

// Format and record the current log.
func (o *log) record(level Level, message string) {
	entity := o.core.getEntity(o, level, message, o.getCaller(level))
	defer o.core.putEntity(entity)

	if err := o.format(entity); err != nil {
		// When the format log fails, we terminate the logging and report the error.
		internal.EchoError("(%s) Failed to format log: %s", o.core.name, err)
	} else {
		err = o.core.hooks.Fire(entity)
		if err != nil {
			internal.EchoError("(%s) Failed to fire log hook: %s", o.core.name, err)
		}
		err = o.write(entity)
		if err != nil {
			internal.EchoError("(%s) Failed to write log: %s", o.core.name, err)
		}
	}

	if level < ErrorLevel {
		switch level {
		case FatalLevel:
			o.core.exitFunc(1)
		case PanicLevel:
			o.core.panicFunc(message)
		}
	}
}

// Format the current log.
func (o *log) format(entity *logEntity) error {
	if formatter := o.core.formatter; formatter != nil {
		return formatter.Format(entity, &entity.buffer)
	}

	kv := map[string]interface{}{
		"name":    entity.name,
		"time":    entity.TimeString(),
		"level":   entity.level.String(),
		"message": entity.message,
	}
	if len(o.fields) > 0 {
		kv["fields"] = internal.StandardiseFieldsForJSONEncoder(o.fields)
	}
	if entity.caller != "" {
		kv["caller"] = entity.caller
	}
	return json.NewEncoder(&entity.buffer).Encode(kv)
}

// Write the current log.
func (o *log) write(entity *logEntity) (err error) {
	var w io.Writer
	if writer, found := o.core.levelWriter[entity.level]; found && writer != nil {
		w = writer
	} else {
		w = o.core.writer
	}
	if o.core.interceptor == nil {
		// When there is no interceptor, make sure that the log written is not empty.
		if entity.Size() > 0 {
			_, err = w.Write(entity.Bytes())
		}
	} else {
		_, err = o.core.interceptor(entity, w)
	}
	return
}

// Get the caller report. If caller reporting is not enabled in the current
// log, an empty string is always returned.
func (o *log) getCaller(level Level) string {
	if o.caller != nil {
		return o.caller.GetCaller()
	}
	if caller, found := o.core.levelCaller[level]; found {
		return caller.GetCaller()
	}
	if o.core.caller != nil {
		return o.core.caller.GetCaller()
	}
	return ""
}

// Log uses the given parameters to record a log of the specified level.
// If the given log level is PanicLevel, the given panic function will be
// called automatically after logging is completed.
// If the given log level is FatalLevel, the given exit function will be
// called automatically after logging is completed.
// If the given log level is invalid, the log will be discarded.
func (o *log) Log(level Level, args ...interface{}) {
	o.log(level, fmt.Sprint(args...))
}

// Uses the given parameters to record a log of the specified level.
func (o *log) log(level Level, args ...interface{}) {
	if !o.core.level.IsEnabled(level) {
		return
	}
	o.record(level, fmt.Sprint(args...))
}

// Logln uses the given parameters to record a log of the specified level.
// If the given log level is PanicLevel, the given panic function will be
// called automatically after logging is completed.
// If the given log level is FatalLevel, the given exit function will be
// called automatically after logging is completed.
// If the given log level is invalid, the log will be discarded.
func (o *log) Logln(level Level, args ...interface{}) {
	o.logln(level, args...)
}

// Uses the given parameters to record a log of the specified level.
func (o *log) logln(level Level, args ...interface{}) {
	if !o.core.level.IsEnabled(level) {
		return
	}
	s := fmt.Sprintln(args...)
	o.record(level, s[:len(s)-1])
}

// Logf uses the given parameters to record a log of the specified level.
// If the given log level is PanicLevel, the given panic function will be
// called automatically after logging is completed.
// If the given log level is FatalLevel, the given exit function will be
// called automatically after logging is completed.
// If the given log level is invalid, the log will be discarded.
func (o *log) Logf(level Level, format string, args ...interface{}) {
	o.logf(level, format, args...)
}

// Uses the given parameters to record a log of the specified level.
func (o *log) logf(level Level, format string, args ...interface{}) {
	if !o.core.level.IsEnabled(level) {
		return
	}
	o.record(level, fmt.Sprintf(format, args...))
}

// Trace uses the given parameters to record a TraceLevel log.
func (o *log) Trace(args ...interface{}) {
	o.log(TraceLevel, args...)
}

// Traceln uses the given parameters to record a TraceLevel log.
func (o *log) Traceln(args ...interface{}) {
	o.logln(TraceLevel, args...)
}

// Tracef uses the given parameters to record a TraceLevel log.
func (o *log) Tracef(format string, args ...interface{}) {
	o.logf(TraceLevel, format, args...)
}

// Print uses the given parameters to record a TraceLevel log.
func (o *log) Print(args ...interface{}) {
	o.log(TraceLevel, args...)
}

// Println uses the given parameters to record a TraceLevel log.
func (o *log) Println(args ...interface{}) {
	o.logln(TraceLevel, args...)
}

// Printf uses the given parameters to record a TraceLevel log.
func (o *log) Printf(format string, args ...interface{}) {
	o.logf(TraceLevel, format, args...)
}

// Debug uses the given parameters to record a DebugLevel log.
func (o *log) Debug(args ...interface{}) {
	o.log(DebugLevel, args...)
}

// Debugln uses the given parameters to record a DebugLevel log.
func (o *log) Debugln(args ...interface{}) {
	o.logln(DebugLevel, args...)
}

// Debugf uses the given parameters to record a DebugLevel log.
func (o *log) Debugf(format string, args ...interface{}) {
	o.logf(DebugLevel, format, args...)
}

// Info uses the given parameters to record a InfoLevel log.
func (o *log) Info(args ...interface{}) {
	o.log(InfoLevel, args...)
}

// Infoln uses the given parameters to record a InfoLevel log.
func (o *log) Infoln(args ...interface{}) {
	o.logln(InfoLevel, args...)
}

// Infof uses the given parameters to record a InfoLevel log.
func (o *log) Infof(format string, args ...interface{}) {
	o.logf(InfoLevel, format, args...)
}

// Echo uses the given parameters to record a InfoLevel log.
func (o *log) Echo(args ...interface{}) {
	o.log(InfoLevel, args...)
}

// Echoln uses the given parameters to record a InfoLevel log.
func (o *log) Echoln(args ...interface{}) {
	o.logln(InfoLevel, args...)
}

// Echof uses the given parameters to record a InfoLevel log.
func (o *log) Echof(format string, args ...interface{}) {
	o.logf(InfoLevel, format, args...)
}

// Warn uses the given parameters to record a WarnLevel log.
func (o *log) Warn(args ...interface{}) {
	o.log(WarnLevel, args...)
}

// Warnln uses the given parameters to record a WarnLevel log.
func (o *log) Warnln(args ...interface{}) {
	o.logln(WarnLevel, args...)
}

// Warnf uses the given parameters to record a WarnLevel log.
func (o *log) Warnf(format string, args ...interface{}) {
	o.logf(WarnLevel, format, args...)
}

// Warning uses the given parameters to record a WarnLevel log.
func (o *log) Warning(args ...interface{}) {
	o.log(WarnLevel, args...)
}

// Warningln uses the given parameters to record a WarnLevel log.
func (o *log) Warningln(args ...interface{}) {
	o.logln(WarnLevel, args...)
}

// Warningf uses the given parameters to record a WarnLevel log.
func (o *log) Warningf(format string, args ...interface{}) {
	o.logf(WarnLevel, format, args...)
}

// Error uses the given parameters to record a ErrorLevel log.
func (o *log) Error(args ...interface{}) {
	o.log(ErrorLevel, args...)
}

// Errorln uses the given parameters to record a ErrorLevel log.
func (o *log) Errorln(args ...interface{}) {
	o.logln(ErrorLevel, args...)
}

// Errorf uses the given parameters to record a ErrorLevel log.
func (o *log) Errorf(format string, args ...interface{}) {
	o.logf(ErrorLevel, format, args...)
}

// Fatal uses the given parameters to record a FatalLevel log.
// After the log record is completed, the system will automatically call
// the exit function given in advance.
func (o *log) Fatal(args ...interface{}) {
	o.log(FatalLevel, args...)
}

// Fatalln uses the given parameters to record a FatalLevel log.
// After the log record is completed, the system will automatically call
// the exit function given in advance.
func (o *log) Fatalln(args ...interface{}) {
	o.logln(FatalLevel, args...)
}

// Fatalf uses the given parameters to record a FatalLevel log.
// After the log record is completed, the system will automatically call
// the exit function given in advance.
func (o *log) Fatalf(format string, args ...interface{}) {
	o.logf(FatalLevel, format, args...)
}

// Panic uses the given parameters to record a PanicLevel log.
// After the log record is completed, the system will automatically call
// the panic function given in advance.
func (o *log) Panic(args ...interface{}) {
	o.log(PanicLevel, args...)
}

// Panicln uses the given parameters to record a PanicLevel log.
// After the log record is completed, the system will automatically call
// the panic function given in advance.
func (o *log) Panicln(args ...interface{}) {
	o.logln(PanicLevel, args...)
}

// Panicf uses the given parameters to record a PanicLevel log.
// After the log record is completed, the system will automatically call
// the panic function given in advance.
func (o *log) Panicf(format string, args ...interface{}) {
	o.logf(PanicLevel, format, args...)
}
