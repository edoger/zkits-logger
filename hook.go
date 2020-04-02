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

// Hook is the log hook type.
// Before the log is written, the corresponding level of hooks will be executed.
// The hooks will be executed sequentially in synchronization until an error is
// encountered or all execution is complete.
type Hook func(Summary) error

// Hooks for each log level.
type hooks map[Level][]Hook

// Add hook for given log level.
func (h hooks) add(level Level, hook Hook) {
    h[level] = append(h[level], hook)
}

// Execute hooks of given log level.
func (h hooks) exec(level Level, s Summary) error {
    for _, hook := range h[level] {
        if err := hook(s); err != nil {
            return err
        }
    }
    return nil
}
