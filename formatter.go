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
	"encoding/json"
	"fmt"
	"time"
)

type Formatter interface {
	Format(Entity, *bytes.Buffer) error
}

func NewJSONFormatter() Formatter {
	return &jsonFormatter{}
}

type jsonFormatter struct {
	timeLayout string
}

func (f *jsonFormatter) Format(entity Entity, buffer *bytes.Buffer) error {
	return json.NewEncoder(buffer).Encode(map[string]interface{}{
		"name":    entity.Name(),
		"time":    entity.Time().Format(f.timeLayout),
		"level":   entity.Level().String(),
		"message": entity.Message(),
		"fields":  f.formatFields(entity.Fields()),
	})
}

func (f *jsonFormatter) formatFields(fields map[string]interface{}) map[string]interface{} {
	if len(fields) == 0 {
		return map[string]interface{}{}
	}
	r := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		if err, ok := v.(error); ok {
			r[k] = err.Error()
		} else {
			r[k] = v
		}
	}
	return r
}

func NewTextFormatter() Formatter {
	return &jsonFormatter{}
}

type textFormatter struct {
	timeLayout string
}

func (f *textFormatter) Format(entity Entity, buffer *bytes.Buffer) error {
	buffer.WriteString("[" + entity.Name() + "]")
	buffer.WriteString("[" + entity.Time().Format(f.timeLayout) + "]")
	buffer.WriteString("[" + entity.Level().String() + "] ")
	buffer.WriteString(entity.Message())

	if fields := entity.Fields(); len(fields) > 0 {
		for k, v := range fields {
			buffer.WriteString(" " + k + "=" + f.toString(v))
		}
	}
	buffer.WriteByte('\n')
	return nil
}

func (f *textFormatter) toString(value interface{}) string {
	if value == nil {
		return "<nil>"
	}

	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case error:
		return v.Error()
	case time.Time:
		return v.Format(f.timeLayout)
	case fmt.Stringer:
		return v.String()
	}

	return fmt.Sprint(value)
}
