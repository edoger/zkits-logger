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
	"encoding/json"
	"fmt"
	"sync"

	"github.com/edoger/zkits-logger/internal"
)

// The default json formatter.
var defaultJSONFormatter = MustNewJSONFormatter(nil, false)

// DefaultJSONFormatter returns the default json formatter.
func DefaultJSONFormatter() Formatter {
	return defaultJSONFormatter
}

// NewJSONFormatter creates and returns an instance of the log json formatter.
// The keys parameter is used to modify the default json field name.
// If the full parameter is true, it will always ensure that all fields exist in
// the top-level json object.
func NewJSONFormatter(keys map[string]string, full bool) (Formatter, error) {
	m := map[string]string{
		"name": "name", "time": "time", "level": "level", "message": "message",
		"fields": "fields", "caller": "caller", "stack": "stack",
	}

	structure := true
	if len(keys) > 0 {
		for key, value := range keys {
			if m[key] == "" {
				return nil, fmt.Errorf("invalid json formatter key %q", key)
			}
			// We ignore the case where all fields are mapped as empty, which is more practical.
			if value != "" && m[key] != value {
				structure = false
				m[key] = value
			}
		}
	}
	f := &jsonFormatter{
		name: m["name"], time: m["time"], level: m["level"], message: m["message"],
		fields: m["fields"], caller: m["caller"], stack: m["stack"],
		full: full, structure: structure,
	}
	return f, nil
}

// MustNewJSONFormatter is like NewJSONFormatter, but triggers a panic when an error occurs.
func MustNewJSONFormatter(keys map[string]string, full bool) Formatter {
	f, err := NewJSONFormatter(keys, full)
	if err != nil {
		panic(err)
	}
	return f
}

// The built-in json formatter.
type jsonFormatter struct {
	name      string
	time      string
	level     string
	message   string
	fields    string
	caller    string
	stack     string
	full      bool
	structure bool
}

// Special built-in structure for json serialization.
// The order of fields cannot be changed.
type jsonFormatterObject struct {
	Caller  *string     `json:"caller,omitempty"`
	Fields  interface{} `json:"fields,omitempty"` // map[string]interface{} or struct{}
	Level   string      `json:"level"`
	Message string      `json:"message"`
	Name    string      `json:"name,omitempty"`
	Stack   []string    `json:"stack,omitempty"`
	Time    *string     `json:"time,omitempty"`
}

// The internal temporary object pool of the json formatter.
var jsonFormatterObjectPool = sync.Pool{
	New: func() interface{} {
		return new(jsonFormatterObject)
	},
}

// Obtain a temporary object for the json formatter from the object pool.
func getJSONFormatterObject() *jsonFormatterObject {
	return jsonFormatterObjectPool.Get().(*jsonFormatterObject)
}

// Put the temporary object of the json formatter back into the object pool.
func putJSONFormatterObject(o *jsonFormatterObject) {
	o.Caller = nil
	o.Fields = nil
	o.Level = ""
	o.Message = ""
	o.Name = ""
	o.Stack = nil
	o.Time = nil

	jsonFormatterObjectPool.Put(o)
}

// Format formats the given log entity into character data and writes it to the given buffer.
func (f *jsonFormatter) Format(e Entity, b *bytes.Buffer) error {
	// In most cases, the performance of json serialization of structure is higher than
	// that of json serialization of map. When the json field name has not changed, we
	// try to use structure for json serialization.
	if f.structure {
		o := getJSONFormatterObject()
		defer putJSONFormatterObject(o)

		o.Level = e.Level().String()
		o.Message = e.Message()
		o.Name = e.Name()
		if tm := e.TimeString(); f.full || tm != "" {
			o.Time = &tm
		}
		if fields := e.Fields(); len(fields) > 0 {
			o.Fields = internal.StandardiseFieldsForJSONEncoder(fields)
		} else {
			if f.full {
				o.Fields = struct{}{}
			}
		}
		if caller := e.Caller(); f.full || caller != "" {
			o.Caller = &caller
		}
		if stack := e.Stack(); len(stack) > 0 {
			o.Stack = stack
		} else {
			if f.full {
				o.Stack = []string{}
			}
		}
		// The json.Encoder.Encode method automatically adds line breaks.
		return json.NewEncoder(b).Encode(o)
	}

	// When the json field cannot be predicted in advance, we use map to package the log data.
	// Is there a better solution to improve the efficiency of json serialization?
	kv := map[string]interface{}{
		f.level:   e.Level().String(),
		f.message: e.Message(),
	}
	if name := e.Name(); f.full || name != "" {
		kv[f.name] = name
	}
	if tm := e.TimeString(); f.full || tm != "" {
		kv[f.time] = tm
	}
	if fields := e.Fields(); len(fields) > 0 {
		kv[f.fields] = internal.StandardiseFieldsForJSONEncoder(fields)
	} else {
		if f.full {
			kv[f.fields] = struct{}{}
		}
	}
	if caller := e.Caller(); f.full || caller != "" {
		kv[f.caller] = caller
	}
	if stack := e.Stack(); len(stack) > 0 {
		kv[f.stack] = stack
	} else {
		if f.full {
			kv[f.stack] = []string{}
		}
	}
	// The json.Encoder.Encode method automatically adds line breaks.
	return json.NewEncoder(b).Encode(kv)
}
