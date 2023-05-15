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
	"io"
)

// Formatter interface defines a standard log formatter.
// The log formatter is not necessary for built-in loggers. We serialize logs to
// JSON format by default. This interface only provides a simple way to change
// this default behavior.
type Formatter interface {
	// Format formats the given log entity into character data and writes it to
	// the given buffer. If the error returned is not empty, the log will be discarded
	// and the registered log hook will not be triggered. We will output the error
	// information to os.Stderr.
	Format(Entity, *bytes.Buffer) error
}

// FormatterFunc type defines a log formatter in the form of a function.
type FormatterFunc func(Entity, *bytes.Buffer) error

// Format formats the given log entity into character data and writes it to
// the given buffer.
func (f FormatterFunc) Format(e Entity, b *bytes.Buffer) error {
	return f(e, b)
}

// UnimplementedFormatter defines an empty, unimplemented log formatter.
// This is usually used to bypass log formatting and implement custom loggers with interceptors.
type UnimplementedFormatter struct{}

// Format does nothing here and always returns a nil error.
func (*UnimplementedFormatter) Format(_ Entity, _ *bytes.Buffer) error {
	return nil
}

// FormatOutput defines the format output of the logger.
type FormatOutput interface {
	// Format formats the given log entity and returns the writer to which the log needs to be written.
	// If the error returned is not empty, the log will be discarded and the registered log hook will
	// not be triggered. We will output the error information to os.Stderr.
	// If the returned log writer is nil, the default log writer is used.
	Format(e Entity, b *bytes.Buffer) (io.Writer, error)
}

// FormatOutputFunc type defines a log format output in the form of a function.
type FormatOutputFunc func(e Entity, b *bytes.Buffer) (io.Writer, error)

// Format formats the given log entity and returns the writer to which the log needs to be written.
func (f FormatOutputFunc) Format(e Entity, b *bytes.Buffer) (io.Writer, error) {
	return f(e, b)
}

// NewFormatOutput creates a log format output instance from the given formatter and writer.
func NewFormatOutput(f Formatter, w io.Writer) FormatOutput {
	return &formatOutput{f, w}
}

// This is the built-in log format output wrapper.
type formatOutput struct {
	f Formatter
	w io.Writer
}

// Format formats the given log entity and returns the writer to which the log needs to be written.
func (w *formatOutput) Format(e Entity, b *bytes.Buffer) (io.Writer, error) {
	return w.w, w.f.Format(e, b)
}

// NewLevelPriorityFormatOutput creates a log level priority format output instance from the given parameters.
func NewLevelPriorityFormatOutput(high, low FormatOutput) FormatOutput {
	return &levelPriorityFormatOutput{high, low}
}

// This is the built-in log level priority format output wrapper.
type levelPriorityFormatOutput struct {
	high FormatOutput
	low  FormatOutput
}

// Format formats the given log entity and returns the writer to which the log needs to be written.
func (w *levelPriorityFormatOutput) Format(e Entity, b *bytes.Buffer) (io.Writer, error) {
	if IsHighPriorityLevel(e.Level()) {
		return w.high.Format(e, b)
	}
	return w.low.Format(e, b)
}
