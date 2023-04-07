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

package internal

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Fields type defines the dynamic field collection of the log.
// After Fields are created, their stored keys will not change.
type Fields map[string]interface{}

// MakeFields creates and returns Fields from a given parameter.
func MakeFields(src map[string]interface{}) Fields {
	r := make(Fields, len(src))
	for k, v := range src {
		r[k] = v
	}
	return r
}

// Clone returns a cloned Fields.
// If n is given, the returned fields will be pre-expanded with equal capacity.
func (fs Fields) Clone(n int) Fields {
	if len(fs) == 0 {
		return make(Fields, n)
	}
	r := make(Fields, len(fs)+n)
	for k, v := range fs {
		r[k] = v
	}
	return r
}

// With returns a cloned Fields and adds the given data to it.
func (fs Fields) With(src map[string]interface{}) Fields {
	if len(src) == 0 {
		return fs.Clone(0)
	}
	r := fs.Clone(len(src))
	for k, v := range src {
		r[k] = v
	}
	return r
}

// StandardiseFieldsForJSONEncoder standardizes the given log fields.
func StandardiseFieldsForJSONEncoder(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		switch o := v.(type) {
		case error:
			// The json.Marshal will convert some errors into "{}", we need to call
			// the error.Error() method before JSON encoding.
			dst[k] = errorToString(o)
		default:
			dst[k] = v
		}
	}
	return dst
}

// FormatFieldsToText standardizes the given log fields.
func FormatFieldsToText(src map[string]interface{}) string {
	texts := make([]string, 0, len(src))
	for k, v := range src {
		texts = append(texts, k+"="+ToString(v))
	}
	// Ensure that the order of log extension fields is consistent.
	if len(texts) > 1 {
		sort.Strings(texts)
	}
	return strings.Join(texts, ", ")
}

// FormatPairsToFields standardizes the given pairs to fields.
func FormatPairsToFields(pairs []interface{}) map[string]interface{} {
	fields := make(map[string]interface{}, len(pairs)/2)
	for i, j := 0, len(pairs); i < j; i += 2 {
		if i+1 < j {
			fields[ToString(pairs[i])] = pairs[i+1]
		} else {
			// Can't be the last key-value pair?
			// We tried setting the value to an empty string, but that shouldn't happen.
			fields[ToString(pairs[i])] = ""
		}
	}
	return fields
}

var (
	errorType    = reflect.TypeOf((*error)(nil)).Elem()
	stringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
)

// ToString tries to convert the given variable into a string.
func ToString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return stringerToString(v)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case bool:
		return strconv.FormatBool(v)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'g', -1, 32)
	case []byte:
		return string(v)
	case error:
		return errorToString(v)
	}
	// Not a common type? We try to use reflection for fast conversion to string.
	rv := reflect.ValueOf(value)
	for k := rv.Kind(); k == reflect.Ptr || k == reflect.Interface; k = rv.Kind() {
		if rv.IsNil() {
			return ""
		}
		if rv.Type().AssignableTo(errorType) {
			return errorToString(rv.Interface().(error))
		}
		if rv.Type().AssignableTo(stringerType) {
			return stringerToString(rv.Interface().(fmt.Stringer))
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.String:
		return rv.String()
	case reflect.Int64, reflect.Int, reflect.Int32, reflect.Int16, reflect.Int8:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint64, reflect.Uint, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return string(rv.Bytes())
		}
	}
	// Ultimately, we can only hope that this returns the desired string.
	return fmt.Sprint(value)
}

// Safely convert the given error to a string.
func errorToString(err error) (s string) {
	if err == nil {
		return
	}
	defer func() {
		if v := recover(); v != nil {
			s = "!!PANIC(error.Error)"
		}
	}()
	s = err.Error()
	return
}

// Safely convert the given fmt.Stringer to a string.
func stringerToString(sr fmt.Stringer) (s string) {
	if sr == nil {
		return
	}
	defer func() {
		if v := recover(); v != nil {
			s = "!!PANIC(fmt.Stringer.String)"
		}
	}()
	s = sr.String()
	return
}
