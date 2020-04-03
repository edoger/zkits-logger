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
	"time"
)

// Summary of the log.
type Summary interface {
	Name() string                     // Get the name of the logger.
	Bytes() []byte                    // Get the log content and return it as bytes.
	String() string                   // Get the log content and return it as string.
	Size() int                        // Get the size of the log.
	Level() Level                     // Get the level of the log.
	Time() time.Time                  // Get the creation time of the log.
	Message() string                  // Get the message of the log.
	Fields() map[string]interface{}   // Get the bound data of the log.
	Field(string) (interface{}, bool) // Query the bound data of the log.
	Context() context.Context         // Get the context of the log.
}

// Implementation of Summary interface.
type summary struct {
	name    string
	level   Level
	message string
	time    time.Time
	buffer  *bytes.Buffer
	fields  map[string]interface{}
	ctx     context.Context
}

// Get the name of the logger.
func (o *summary) Name() string {
	return o.name
}

// Get the log content and return it as bytes.
func (o *summary) Bytes() []byte {
	return o.buffer.Bytes()
}

// Get the log content and return it as string.
func (o *summary) String() string {
	return o.buffer.String()
}

// Get the size of the log.
func (o *summary) Size() int {
	return o.buffer.Len()
}

// Get the level of the log.
func (o *summary) Level() Level {
	return o.level
}

// Get the creation time of the log.
func (o *summary) Time() time.Time {
	return o.time
}

// Get the message of the log.
func (o *summary) Message() string {
	return o.message
}

// Get the bound data of the log.
func (o *summary) Fields() map[string]interface{} {
	return o.fields
}

// Query the bound data of the log.
func (o *summary) Field(key string) (value interface{}, found bool) {
	if len(o.fields) > 0 {
		value, found = o.fields[key]
	}
	return
}

// Get the context of the log.
func (o *summary) Context() context.Context {
	if o.ctx == nil {
		return context.Background()
	}
	return o.ctx
}
