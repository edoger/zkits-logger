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

type Entity interface {
	// Name returns the logger name.
	Name() string

	// Time returns the log time.
	Time() time.Time

	// Level returns the log level.
	Level() Level

	// Message returns the log message.
	Message() string

	// Fields returns the log fields.
	Fields() map[string]interface{}

	// Context returns the log context.
	Context() context.Context

	// Caller returns the log caller.
	// If it is not enabled, an empty string is always returned.
	Caller() string
}

// Summary interface defines the summary of the log. The log summary is the final state of a log,
// and the content of the log will no longer change.
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
}

// The logEntity type is a built-in implementation of the Entity interface.
type logEntity struct {
	name    string
	time    time.Time
	level   Level
	message string
	fields  map[string]interface{}
	ctx     context.Context
	buffer  bytes.Buffer
	caller  string
}

// Name returns the logger name.
func (o *logEntity) Name() string {
	return o.name
}

// Time returns the log time.
func (o *logEntity) Time() time.Time {
	return o.time
}

// Level returns the log level.
func (o *logEntity) Level() Level {
	return o.level
}

// Message returns the log message.
func (o *logEntity) Message() string {
	return o.message
}

// Fields returns the log fields.
func (o *logEntity) Fields() map[string]interface{} {
	return o.fields
}

// Context returns the log context.
func (o *logEntity) Context() context.Context {
	if o.ctx == nil {
		return context.Background()
	}
	return o.ctx
}

// Caller returns the log caller.
// If it is not enabled, an empty string is always returned.
func (o *logEntity) Caller() string {
	return o.caller
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

// Read is the implementation of io.Reader interface.
func (o *logEntity) Read(p []byte) (int, error) {
	return o.buffer.Read(p)
}
