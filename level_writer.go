// Copyright 2022 The ZKits Project Authors.
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
)

// NewLevelWriter creates a writer that records each message written as a log message.
// Generally, this writer is used to bridge the loggers of the standard library.
func NewLevelWriter(l Level, g Log) io.Writer {
	return &logLevelWriter{level: l, log: g}
}

// This writer records each message written as a log message.
// The level of the log is determined by the level given when the writer was created.
type logLevelWriter struct {
	level Level
	log   Log
}

// Write method is an implementation of the io.Writer interface.
// This method will always write successfully.
func (w *logLevelWriter) Write(p []byte) (n int, err error) {
	// Strips extra newlines, as the log formatter automatically appends a newline.
	if n = len(p); n > 0 && p[n-1] == '\n' {
		w.log.Log(w.level, string(p[:n-1]))
	} else {
		w.log.Log(w.level, string(p))
	}
	return
}
