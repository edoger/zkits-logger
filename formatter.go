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
