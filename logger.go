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
)

// Definition of the logger.
type Logger interface {
	Log

	// Get the current logger level.
	GetLevel() Level
	// Set the current logger level.
	SetLevel(Level) Logger
	// Get the current logger output writer.
	GetOutput() io.Writer
	// Set the current logger output writer.
	SetOutput(io.Writer) Logger
	// Set the current logger system exit function.
	WithExitFunc(func(int)) Logger
	// Register hook to the current logger.
	WithHook([]Level, Hook) Logger
}

// Create a new logger.
// By default, the logger level is InfoLevel and logs will be output to
// standard output.
func New(name string) Logger {
	return &logger{
		log: log{
			common: newCommon(name),
			fields: make(map[string]interface{}),
		},
	}
}

// Internal implementation of the Logger interface.
type logger struct {
	log
}

// Get the current logger level.
func (o *logger) GetLevel() Level {
	return o.log.common.level
}

// Set the current logger level.
func (o *logger) SetLevel(level Level) Logger {
	o.log.common.mutex.Lock()
	defer o.log.common.mutex.Unlock()
	// Invalid level will be ignored.
	if level.IsValid() {
		o.log.common.level = level
	}
	return o
}

// Get the current logger output writer.
func (o *logger) GetOutput() io.Writer {
	return o.log.common.writer
}

// Set the current logger output writer.
func (o *logger) SetOutput(w io.Writer) Logger {
	o.log.common.mutex.Lock()
	defer o.log.common.mutex.Unlock()

	o.log.common.writer = w
	return o
}

// Set the current logger system exit function.
func (o *logger) WithExitFunc(exit func(int)) Logger {
	o.log.common.mutex.Lock()
	defer o.log.common.mutex.Unlock()

	if exit == nil {
		o.log.common.exit = os.Exit
	} else {
		o.log.common.exit = exit
	}
	return o
}

// Register hook to the current logger.
func (o *logger) WithHook(levels []Level, hook Hook) Logger {
	o.log.common.mutex.Lock()
	defer o.log.common.mutex.Unlock()

	for _, level := range levels {
		o.log.common.hooks.add(level, hook)
	}
	return o
}
