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
	"io"
	"time"
)

// Entity interface defines the entity of the log.
type Entity interface {
	// Name returns the logger name.
	Name() string

	// Time returns the log time.
	Time() time.Time

	// TimeString returns the log time string formatted with the default time format.
	TimeString() string

	// Level returns the log level.
	Level() Level

	// Message returns the log message.
	Message() string

	// HasFields determines whether the log contains fields.
	HasFields() bool

	// Fields returns the log fields.
	Fields() map[string]interface{}

	// HasContext determines whether the log contains a context.
	HasContext() bool

	// Context returns the log context.
	Context() context.Context

	// HasCaller determines whether the log contains caller information.
	HasCaller() bool

	// Caller returns the log caller.
	// If it is not enabled, an empty string is always returned.
	Caller() string

	// HasStack determines whether the log contains call stack information.
	HasStack() bool

	// Stack returns the call stack information at the logging location.
	// Returns nil if not enabled.
	Stack() []string

	// Buffer returns the entity buffer instance.
	Buffer() *bytes.Buffer
}

// Summary interface defines the summary of the log.
// The log summary is the final state of a log, and the content of the log will no longer change.
type Summary interface {
	Entity
	io.Reader

	// Bytes returns the log content bytes.
	// This method returns the content processed by the formatter.
	Bytes() []byte

	// String returns the log content string.
	// This method returns the content processed by the formatter.
	String() string

	// Size returns the log content size.
	Size() int

	// Clone returns a copy of the current log summary (excluding context).
	Clone() Summary

	// CloneWithContext returns a copy of the current log summary and sets its context to the given value.
	CloneWithContext(context.Context) Summary
}

// The logEntity type is a built-in implementation of the Entity interface.
type logEntity struct {
	name       string
	time       time.Time
	timeFormat string
	level      Level
	message    string
	fields     map[string]interface{}
	ctx        context.Context
	buffer     bytes.Buffer
	caller     string
	stack      []string
}

// Name returns the logger name.
func (o *logEntity) Name() string {
	return o.name
}

// Time returns the log time.
func (o *logEntity) Time() time.Time {
	return o.time
}

// TimeString returns the log time string formatted with the default time format.
func (o *logEntity) TimeString() string {
	if o.timeFormat == "" {
		return ""
	}
	return o.time.Format(o.timeFormat)
}

// Level returns the log level.
func (o *logEntity) Level() Level {
	return o.level
}

// Message returns the log message.
func (o *logEntity) Message() string {
	return o.message
}

// HasFields determines whether the log contains fields.
func (o *logEntity) HasFields() bool {
	return len(o.fields) > 0
}

// Fields returns the log fields.
func (o *logEntity) Fields() map[string]interface{} {
	return o.fields
}

// HasContext determines whether the log contains a context.
func (o *logEntity) HasContext() bool {
	return o.ctx != nil
}

// Context returns the log context.
func (o *logEntity) Context() context.Context {
	if o.ctx == nil {
		return context.Background()
	}
	return o.ctx
}

// HasCaller determines whether the log contains caller information.
func (o *logEntity) HasCaller() bool {
	return o.caller != ""
}

// Caller returns the log caller.
// If it is not enabled, an empty string is always returned.
func (o *logEntity) Caller() string {
	return o.caller
}

// HasStack determines whether the log contains call stack information.
func (o *logEntity) HasStack() bool {
	return len(o.stack) > 0
}

// Stack returns the call stack information at the logging location.
// Returns nil if not enabled.
func (o *logEntity) Stack() []string {
	return o.stack
}

// Buffer returns the entity buffer instance.
func (o *logEntity) Buffer() *bytes.Buffer {
	return &o.buffer
}

// Bytes returns the log content bytes.
// This method returns the content processed by the formatter.
func (o *logEntity) Bytes() []byte {
	return o.buffer.Bytes()
}

// String returns the log content string.
// This method returns the content processed by the formatter.
func (o *logEntity) String() string {
	return o.buffer.String()
}

// Size returns the log content size.
func (o *logEntity) Size() int {
	return o.buffer.Len()
}

// Clone returns a copy of the current log summary (excluding context).
func (o *logEntity) Clone() Summary {
	return o.CloneWithContext(nil)
}

// CloneWithContext returns a copy of the current log summary and sets its context to the given value.
func (o *logEntity) CloneWithContext(ctx context.Context) Summary {
	var buffer *bytes.Buffer
	if bs := o.buffer.Bytes(); len(bs) > 0 {
		cp := make([]byte, len(bs))
		copy(cp, bs)
		buffer = bytes.NewBuffer(cp)
	} else {
		buffer = new(bytes.Buffer)
	}

	var fields map[string]interface{}
	if n := len(o.fields); n > 0 {
		fields = make(map[string]interface{}, n)
		for k, v := range o.fields {
			fields[k] = v
		}
	}

	var stack []string
	if o.stack != nil {
		stack = make([]string, len(o.stack))
		copy(stack, o.stack)
	}

	return &logEntity{
		name:       o.name,
		time:       o.time,
		timeFormat: o.timeFormat,
		level:      o.level,
		message:    o.message,
		fields:     fields,
		ctx:        ctx,
		buffer:     *buffer,
		caller:     o.caller,
		stack:      stack,
	}
}

// Read is the implementation of io.Reader interface.
func (o *logEntity) Read(p []byte) (int, error) {
	return o.buffer.Read(p)
}
