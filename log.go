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
)

type Log interface {
	// Gets the name of the logger.
	Name() string
	// Add an extra data to the current log.
	// This method will return a new log instance.
	WithField(string, interface{}) Log
	// Add error to the current log.
	// This method will return a new log instance.
	// This method is relative to WithField("error", error).
	WithError(error) Log
	// Add extra data to the current log.
	// This method will return a new log instance.
	WithFields(map[string]interface{}) Log
	// Add context to the current log.
	// This method will return a new log instance.
	WithContext(context.Context) Log
	// Record the current log and specify the level.
	Log(Level, ...interface{})
	Logln(Level, ...interface{})
	Logf(Level, string, ...interface{})
	// Record the current log and specify the level as TraceLevel.
	Trace(...interface{})
	Traceln(...interface{})
	Tracef(string, ...interface{})
	Print(...interface{})
	Println(...interface{})
	Printf(string, ...interface{})
	// Record the current log and specify the level as DebugLevel.
	Debug(...interface{})
	Debugln(...interface{})
	Debugf(string, ...interface{})
	// Record the current log and specify the level as InfoLevel.
	Info(...interface{})
	Infoln(...interface{})
	Infof(string, ...interface{})
	// Record the current log and specify the level as WarnLevel.
	Warn(...interface{})
	Warnln(...interface{})
	Warnf(string, ...interface{})
	Warning(...interface{})
	Warningln(...interface{})
	Warningf(string, ...interface{})
	// Record the current log and specify the level as ErrorLevel.
	Error(...interface{})
	Errorln(...interface{})
	Errorf(string, ...interface{})
	// Record the current log and specify the level as FatalLevel.
	// The exit function will be executed automatically.
	Fatal(...interface{})
	Fatalln(...interface{})
	Fatalf(string, ...interface{})
	// Record the current log and specify the level as PanicLevel.
	// Panic will be triggered automatically.
	Panic(...interface{})
	Panicln(...interface{})
	Panicf(string, ...interface{})
}

// Common properties of logger instance and log instance.
type common struct {
	// Mutex for shared resource operations.
	mutex sync.Mutex
	// The name of the logger.
	name string
	// The level of the logger.
	level Level
	// The system exit function.
	// This function is called after logging at the FatalLevel level.
	// By default, it is os.Exit().
	exit func(int)
	// Registered log hooks.
	hooks hooks
	// The target to which all logs are finally written.
	writer io.Writer
}

// Create a new common instance and bind the logger name.
func newCommon(name string) *common {
	return &common{
		name:   name,
		level:  InfoLevel,
		exit:   os.Exit,
		hooks:  make(hooks),
		writer: os.Stdout,
	}
}

// Gets the name of the logger.
func (c *common) Name() string {
	return c.name
}

// Internal implementation of the Log interface.
type log struct {
	*common
	// Bind data for the current log.
	fields map[string]interface{}
	// Bind context for the current log.
	ctx context.Context
}

// Add an extra data to the current log.
// This method will return a new log instance.
func (o *log) WithField(key string, value interface{}) Log {
	return o.WithFields(map[string]interface{}{key: value})
}

// Add error to the current log.
// This method will return a new log instance.
func (o *log) WithError(err error) Log {
	return o.WithFields(map[string]interface{}{"error": err})
}

// Add extra data to the current log.
// This method will return a new log instance.
func (o *log) WithFields(fields map[string]interface{}) Log {
	r := &log{
		common: o.common,
		fields: make(map[string]interface{}, len(o.fields)+len(fields)),
		ctx:    o.ctx,
	}
	for k, v := range o.fields {
		r.fields[k] = v
	}
	for k, v := range fields {
		r.fields[k] = v
	}
	return r
}

// Add context to the current log.
// This method will return a new log instance.
func (o *log) WithContext(ctx context.Context) Log {
	r := &log{
		common: o.common,
		fields: make(map[string]interface{}, len(o.fields)),
		ctx:    ctx,
	}
	for k, v := range o.fields {
		r.fields[k] = v
	}
	return r
}

// Write the current log to the output destination.
// If necessary, the corresponding level of log hooks are also executed.
// For FatalLevel log, the exit function will be executed automatically.
// For PanicLevel log, panic will be triggered automatically.
func (o *log) log(level Level, message string) {
	fields := make(map[string]interface{}, len(o.fields))
	for k, v := range o.fields {
		switch v := v.(type) {
		case error:
			fields[k] = v.Error()
		default:
			fields[k] = v
		}
	}

	now := time.Now()
	buffer := new(bytes.Buffer)
	// Newline characters are automatically appended when serializing
	// with json.Encoder.
	err := json.NewEncoder(buffer).Encode(map[string]interface{}{
		"time":    now.Format("2006-01-02 15:04:05"),
		"logger":  o.common.name,
		"level":   level.String(),
		"message": message,
		"fields":  fields,
	})
	if err != nil {
		echo("Failed to serialize log: %s", err)
		return
	}

	su := new(summary)
	su.name = o.common.name
	su.level = level
	su.message = message
	su.time = now
	su.buffer = buffer
	su.ctx = o.ctx
	su.fields = make(map[string]interface{}, len(o.fields))
	for k, v := range o.fields {
		su.fields[k] = v
	}
	if err := o.common.hooks.exec(level, su); err != nil {
		echo("Failed to fire hook: %s", err)
	}

	if _, err := o.common.writer.Write(buffer.Bytes()); err != nil {
		echo("Failed to write log: %s", err)
	}

	switch level {
	case FatalLevel:
		o.common.exit(1)
	case PanicLevel:
		panic(su)
	}
}

// Record the current log and specify the level.
func (o *log) Log(level Level, args ...interface{}) {
	if level <= o.common.level {
		o.log(level, fmt.Sprint(args...))
	}
}

// Record the current log and specify the level.
func (o *log) Logln(level Level, args ...interface{}) {
	if level <= o.common.level {
		message := fmt.Sprintln(args...)
		o.log(level, message[:len(message)-1])
	}
}

// Record the current log and specify the level.
func (o *log) Logf(level Level, format string, args ...interface{}) {
	if level <= o.common.level {
		o.log(level, fmt.Sprintf(format, args...))
	}
}

// Record the current log and specify the level as TraceLevel.
func (o *log) Trace(args ...interface{}) {
	o.Log(TraceLevel, args...)
}

// Record the current log and specify the level as TraceLevel.
func (o *log) Traceln(args ...interface{}) {
	o.Logln(TraceLevel, args...)
}

// Record the current log and specify the level as TraceLevel.
func (o *log) Tracef(format string, args ...interface{}) {
	o.Logf(TraceLevel, format, args...)
}

// Record the current log and specify the level as TraceLevel.
func (o *log) Print(args ...interface{}) {
	o.Log(TraceLevel, args...)
}

// Record the current log and specify the level as TraceLevel.
func (o *log) Println(args ...interface{}) {
	o.Logln(TraceLevel, args...)
}

// Record the current log and specify the level as TraceLevel.
func (o *log) Printf(format string, args ...interface{}) {
	o.Logf(TraceLevel, format, args...)
}

// Record the current log and specify the level as DebugLevel.
func (o *log) Debug(args ...interface{}) {
	o.Log(DebugLevel, args...)
}

// Record the current log and specify the level as DebugLevel.
func (o *log) Debugln(args ...interface{}) {
	o.Logln(DebugLevel, args...)
}

// Record the current log and specify the level as DebugLevel.
func (o *log) Debugf(format string, args ...interface{}) {
	o.Logf(DebugLevel, format, args...)
}

// Record the current log and specify the level as InfoLevel.
func (o *log) Info(args ...interface{}) {
	o.Log(InfoLevel, args...)
}

// Record the current log and specify the level as InfoLevel.
func (o *log) Infoln(args ...interface{}) {
	o.Logln(InfoLevel, args...)
}

// Record the current log and specify the level as InfoLevel.
func (o *log) Infof(format string, args ...interface{}) {
	o.Logf(InfoLevel, format, args...)
}

// Record the current log and specify the level as WarnLevel.
func (o *log) Warn(args ...interface{}) {
	o.Log(WarnLevel, args...)
}

// Record the current log and specify the level as WarnLevel.
func (o *log) Warnln(args ...interface{}) {
	o.Logln(WarnLevel, args...)
}

// Record the current log and specify the level as WarnLevel.
func (o *log) Warnf(format string, args ...interface{}) {
	o.Logf(WarnLevel, format, args...)
}

// Record the current log and specify the level as WarnLevel.
func (o *log) Warning(args ...interface{}) {
	o.Log(WarnLevel, args...)
}

// Record the current log and specify the level as WarnLevel.
func (o *log) Warningln(args ...interface{}) {
	o.Logln(WarnLevel, args...)
}

// Record the current log and specify the level as WarnLevel.
func (o *log) Warningf(format string, args ...interface{}) {
	o.Logf(WarnLevel, format, args...)
}

// Record the current log and specify the level as ErrorLevel.
func (o *log) Error(args ...interface{}) {
	o.Log(ErrorLevel, args...)
}

// Record the current log and specify the level as ErrorLevel.
func (o *log) Errorln(args ...interface{}) {
	o.Logln(ErrorLevel, args...)
}

// Record the current log and specify the level as ErrorLevel.
func (o *log) Errorf(format string, args ...interface{}) {
	o.Logf(ErrorLevel, format, args...)
}

// Record the current log and specify the level as FatalLevel.
// The exit function will be executed automatically.
func (o *log) Fatal(args ...interface{}) {
	o.Log(FatalLevel, args...)
}

// Record the current log and specify the level as FatalLevel.
// The exit function will be executed automatically.
func (o *log) Fatalln(args ...interface{}) {
	o.Logln(FatalLevel, args...)
}

// Record the current log and specify the level as FatalLevel.
// The exit function will be executed automatically.
func (o *log) Fatalf(format string, args ...interface{}) {
	o.Logf(FatalLevel, format, args...)
}

// Record the current log and specify the level as PanicLevel.
// Panic will be triggered automatically.
func (o *log) Panic(args ...interface{}) {
	o.Log(PanicLevel, args...)
}

// Record the current log and specify the level as PanicLevel.
// Panic will be triggered automatically.
func (o *log) Panicln(args ...interface{}) {
	o.Logln(PanicLevel, args...)
}

// Record the current log and specify the level as PanicLevel.
// Panic will be triggered automatically.
func (o *log) Panicf(format string, args ...interface{}) {
	o.Logf(PanicLevel, format, args...)
}
