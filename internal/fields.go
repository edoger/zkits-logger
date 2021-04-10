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
	"sort"
	"strings"
)

// Fields type defines the dynamic field collection of the log.
// After Fields are created, their stored keys will not change.
type Fields map[string]interface{}

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
			dst[k] = o.Error()
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
		switch o := v.(type) {
		case []byte:
			texts = append(texts, k+"="+string(o))
		default:
			texts = append(texts, k+"="+fmt.Sprint(v))
		}
	}
	// Ensure that the order of log extension fields is consistent.
	if len(texts) > 1 {
		sort.Strings(texts)
	}
	return strings.Join(texts, ", ")
}
