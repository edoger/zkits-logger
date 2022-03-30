// Copyright 2021 The ZKits Project Authors.
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
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/edoger/zkits-logger/internal"
)

// This regular expression is used to analyze placeholders in text formatter format.
var formatRegexp = regexp.MustCompile(`{(name|time|level|message|caller|stack|fields)(?:@?([^{}]*)?)?}`)

// The default text formatter.
var defaultTextFormatter = MustNewTextFormatter("{name}:[{time}][{level@sc}] {message}{caller}{fields}{stack}", false)

// The default quote text formatter.
var defaultQuoteTextFormatter = MustNewTextFormatter("{name}:[{time}][{level@sc}] {message}{caller}{fields}{stack}", true)

// DefaultTextFormatter returns the default text formatter.
func DefaultTextFormatter() Formatter {
	return defaultTextFormatter
}

// DefaultQuoteTextFormatter returns the default quote text formatter.
func DefaultQuoteTextFormatter() Formatter {
	return defaultQuoteTextFormatter
}

// NewTextFormatter creates and returns an instance of the log text formatter.
// The format parameter is used to control the format of the log, and it has many control parameters.
// For example:
//     "[{name}][{time}][{level}] {caller} {message} {fields}"
// The supported placeholders are:
//     {name}      The name of the logger that recorded this log.
//     {time}      Record the time of this log.
//     {level}     The level of this log.
//     {caller}    The name and line number of the file where this log was generated. (If enabled)
//     {message}   The message of this log.
//     {fields}    The extended fields of this log. (if it exists)
// It is worth knowing:
//     1. For the {time} parameter, we can specify time format, like this: {time@2006-01-02 15:04:05}.
//     2. For the {level} parameter, we can specify level format, like this: {level@sc},
//        {level@sc} or {level@cs} will call the Level.ShortCapitalString method.
//        {level@s} will call the Level.ShortString method.
//        {level@c} will call the Level.CapitalString method.
//        For other will call the Level.String method.
//     3. Considering the aesthetics of the format, for {caller} and {fields}, if
//        there is non-empty data, a space will be automatically added in front.
//        If this behavior is not needed, use {caller@?} or {fields@?} parameters.
// The quote parameter is used to escape invisible characters in the log.
func NewTextFormatter(format string, quote bool) (Formatter, error) {
	sub := formatRegexp.FindAllStringSubmatch(format, -1)
	if len(sub) == 0 {
		return nil, fmt.Errorf("invalid text formatter format %q", format)
	}
	// If sub is not empty, then idx is definitely not empty.
	idx := formatRegexp.FindAllStringIndex(format, -1)
	f := &textFormatter{quote: quote, callerPrefix: " ", fieldsPrefix: " ", stackPrefix: " "}

	var parts []string
	var start int
	for i, j := 0, len(sub); i < j; i++ {
		key, args := sub[i][1], sub[i][2]
		switch key {
		case "name":
			f.encoders = append(f.encoders, f.encodeName)
		case "time":
			if args == "" {
				f.encoders = append(f.encoders, f.encodeTime)
			} else {
				f.timeFormat = args
				f.encoders = append(f.encoders, f.encodeTimeWithFormat)
			}
		case "level":
			switch args {
			case "cs", "sc":
				f.encoders = append(f.encoders, f.encodeShortCapitalLevel)
			case "s":
				f.encoders = append(f.encoders, f.encodeShortLevel)
			case "c":
				f.encoders = append(f.encoders, f.encodeCapitalLevel)
			default:
				f.encoders = append(f.encoders, f.encodeLevel)
			}
		case "message":
			f.encoders = append(f.encoders, f.encodeMessage)
		case "caller":
			f.encoders = append(f.encoders, f.encodeCaller)
			if args == "?" {
				f.callerPrefix = ""
			}
		case "fields":
			f.encoders = append(f.encoders, f.encodeFields)
			if args == "?" {
				f.fieldsPrefix = ""
			}
		case "stack":
			f.encoders = append(f.encoders, f.encodeStack)
			if args == "?" {
				f.stackPrefix = ""
			}
		}
		parts = append(parts, format[start:idx[i][0]])
		start = idx[i][1]
	}

	f.format = strings.Join(append(parts, format[start:]), "%s")
	return f, nil
}

// MustNewTextFormatter is like NewTextFormatter, but triggers a panic when an error occurs.
func MustNewTextFormatter(format string, quote bool) Formatter {
	f, err := NewTextFormatter(format, quote)
	if err != nil {
		panic(err)
	}
	return f
}

// The built-in text formatter.
type textFormatter struct {
	format       string
	quote        bool
	encoders     []func(Entity) string
	timeFormat   string
	callerPrefix string
	fieldsPrefix string
	stackPrefix  string
}

// Format formats the given log entity into character data and writes it to the given buffer.
func (f *textFormatter) Format(e Entity, b *bytes.Buffer) (err error) {
	args := make([]interface{}, len(f.encoders))
	for i, j := 0, len(f.encoders); i < j; i++ {
		args[i] = f.encoders[i](e)
	}
	if f.quote {
		// The quoted[0] and quoted[len(s)-1] is '"', they need to be removed.
		quoted := strconv.AppendQuote(nil, fmt.Sprintf(f.format, args...))
		quoted[len(quoted)-1] = '\n'
		_, err = b.Write(quoted[1:])
	} else {
		_, err = b.WriteString(fmt.Sprintf(f.format, args...) + "\n")
	}
	return
}

// Encode the name of the log.
func (f *textFormatter) encodeName(e Entity) string {
	return e.Name()
}

// Encode the level of the log.
func (f *textFormatter) encodeLevel(e Entity) string {
	return e.Level().String()
}

// Encode the short capital level of the log.
func (f *textFormatter) encodeShortCapitalLevel(e Entity) string {
	return e.Level().ShortCapitalString()
}

// Encode the short level of the log.
func (f *textFormatter) encodeShortLevel(e Entity) string {
	return e.Level().ShortString()
}

// Encode the capital level of the log.
func (f *textFormatter) encodeCapitalLevel(e Entity) string {
	return e.Level().CapitalString()
}

// Encode the time of the log.
func (f *textFormatter) encodeTime(e Entity) string {
	return e.TimeString()
}

// Encode the time with format of the log.
func (f *textFormatter) encodeTimeWithFormat(e Entity) string {
	return e.Time().Format(f.timeFormat)
}

// Encode the caller of the log.
func (f *textFormatter) encodeCaller(e Entity) string {
	if s := e.Caller(); s != "" {
		return f.callerPrefix + s
	}
	return ""
}

// Encode the message of the log.
func (f *textFormatter) encodeMessage(e Entity) string {
	return e.Message()
}

// Encode the fields of the log.
func (f *textFormatter) encodeFields(e Entity) string {
	if fields := e.Fields(); len(fields) > 0 {
		return f.fieldsPrefix + internal.FormatFieldsToText(e.Fields())
	}
	return ""
}

// Encode the stack of the log.
func (f *textFormatter) encodeStack(e Entity) string {
	if stack := e.Stack(); len(stack) > 0 {
		return f.stackPrefix + strings.Join(stack, "; ")
	}
	return ""
}
