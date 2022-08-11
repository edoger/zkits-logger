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

// NewJSONFormatterFromPool creates a JSON formatter from the given JSON serializable object pool.
func NewJSONFormatterFromPool(p JSONFormatterObjectPool) Formatter {
	return &jsonFormatter{pool: p}
}

// NewJSONFormatter creates and returns an instance of the log json formatter.
// The keys parameter is used to modify the default json field name.
// If the full parameter is true, it will always ensure that all fields exist in the top-level json object.
func NewJSONFormatter(keys map[string]string, full bool) (Formatter, error) {
	if len(keys) > 0 {
		structure := true
		mapping := map[string]string{
			"name": "name", "time": "time", "level": "level", "message": "message",
			"fields": "fields", "caller": "caller", "stack": "stack",
		}
		for key, value := range keys {
			if mapping[key] == "" {
				// We require that the key-name map must be pure.
				return nil, fmt.Errorf("invalid json formatter key %q", key)
			}
			// We ignore the case where all fields are mapped as empty, which is more practical.
			if value != "" && mapping[key] != value {
				structure = false
				mapping[key] = value
			}
		}
		// when the json field cannot be predicted in advance, we use map to package the log data.
		// is there a better solution to improve the efficiency of json serialization?
		if !structure {
			return NewJSONFormatterFromPool(newJSONFormatterMapPool(full, mapping)), nil
		}
	}
	// In most cases, the performance of json serialization of structure is higher than
	// that of json serialization of map. When the json field name has not changed, we
	// try to use structure for json serialization.
	return NewJSONFormatterFromPool(newJSONFormatterObjectPool(full)), nil
}

// MustNewJSONFormatter is like NewJSONFormatter, but triggers a panic when an error occurs.
func MustNewJSONFormatter(keys map[string]string, full bool) Formatter {
	f, err := NewJSONFormatter(keys, full)
	if err != nil {
		panic(err)
	}
	return f
}

// JSONFormatterObjectPool defines a pool of serializable objects for JSON formatter.
// This object pool is used to create and recycle json log objects.
type JSONFormatterObjectPool interface {
	// GetObject creates and returns a new JSON log object from the given log Entity.
	// The returned object must be JSON-Serializable, otherwise the formatter will fail to work.
	GetObject(Entity) interface{}

	// PutObject recycles json serialized objects.
	// This method must clean up the log data bound to the object to free memory.
	PutObject(interface{})
}

// The built-in json formatter.
type jsonFormatter struct {
	pool JSONFormatterObjectPool
}

// Format formats the given log entity into character data and writes it to the given buffer.
func (f *jsonFormatter) Format(e Entity, b *bytes.Buffer) (err error) {
	o := f.pool.GetObject(e)
	// The json.Encoder.Encode method automatically adds line breaks.
	err = json.NewEncoder(b).Encode(o)
	f.pool.PutObject(o)
	return
}

// This is the built-in pool of serializable JSON map.
type jsonFormatterMapPool struct {
	full bool
	// These fields store the names of the keys in the json object.
	name, time, level, message, fields, caller, stack string
}

// Creates and returns a new pool of serializable JSON map.
func newJSONFormatterMapPool(full bool, keys map[string]string) JSONFormatterObjectPool {
	return &jsonFormatterMapPool{
		full: full, name: keys["name"], time: keys["time"], level: keys["level"],
		message: keys["message"], fields: keys["fields"], caller: keys["caller"], stack: keys["stack"],
	}
}

// GetObject creates and returns a new JSON log map from the given log Entity.
// This method is an implementation of the JSONFormatterObjectPool interface.
func (p *jsonFormatterMapPool) GetObject(e Entity) interface{} {
	kv := map[string]interface{}{p.level: e.Level().String(), p.message: e.Message()}
	if name := e.Name(); p.full || name != "" {
		kv[p.name] = name
	}
	if tm := e.TimeString(); p.full || tm != "" {
		kv[p.time] = tm
	}
	if fields := e.Fields(); len(fields) > 0 {
		kv[p.fields] = internal.StandardiseFieldsForJSONEncoder(fields)
	} else {
		if p.full { // Always keep it as an empty json object.
			kv[p.fields] = struct{}{}
		}
	}
	if caller := e.Caller(); p.full || caller != "" {
		kv[p.caller] = caller
	}
	if stack := e.Stack(); len(stack) > 0 {
		kv[p.stack] = stack
	} else {
		if p.full {
			kv[p.stack] = []string{}
		}
	}
	return kv
}

// PutObject does nothing here.
// This method is an implementation of the JSONFormatterObjectPool interface.
func (*jsonFormatterMapPool) PutObject(interface{}) { /* do nothing */ }

// This is the built-in pool of serializable JSON objects.
type jsonFormatterObjectPool struct {
	full bool
	pool *sync.Pool
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

// Creates and returns a new pool of serializable JSON objects.
func newJSONFormatterObjectPool(full bool) JSONFormatterObjectPool {
	return &jsonFormatterObjectPool{full: full, pool: &sync.Pool{
		New: func() interface{} { return new(jsonFormatterObject) },
	}}
}

// GetObject creates and returns a new JSON log object from the given log Entity.
// This method is an implementation of the JSONFormatterObjectPool interface.
func (p *jsonFormatterObjectPool) GetObject(e Entity) interface{} {
	o := p.pool.Get().(*jsonFormatterObject)
	o.Level, o.Message, o.Name = e.Level().String(), e.Message(), e.Name()
	if tm := e.TimeString(); p.full || tm != "" {
		o.Time = &tm
	}
	if fields := e.Fields(); len(fields) > 0 {
		o.Fields = internal.StandardiseFieldsForJSONEncoder(fields)
	} else {
		if p.full { // Always keep it as an empty json object.
			o.Fields = struct{}{}
		}
	}
	if caller := e.Caller(); p.full || caller != "" {
		o.Caller = &caller
	}
	if stack := e.Stack(); len(stack) > 0 {
		o.Stack = stack
	} else {
		if p.full {
			o.Stack = []string{}
		}
	}
	return o
}

// PutObject recycles json serialized objects.
// This method is an implementation of the JSONFormatterObjectPool interface.
func (p *jsonFormatterObjectPool) PutObject(v interface{}) {
	o := v.(*jsonFormatterObject)
	o.Caller, o.Fields, o.Level, o.Message, o.Name, o.Stack, o.Time = nil, nil, "", "", "", nil, nil
	p.pool.Put(o)
}
