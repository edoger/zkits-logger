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

// Hook interface defines the log hook program.
// The log hook is triggered synchronously after the log is successfully formatted and
// before it is written.
// Any error returned by the hook program will not prevent the log from being written.
type Hook interface {
	// Levels returns the log levels associated with the current log hook.
	// This method is only called once when the log hook is added to the logger.
	// If the returned log levels are empty, the hook will not be added.
	// If there are invalid log level in the returned log levels, it will be automatically ignored.
	Levels() []Level

	// Fire receives the summary of the log and performs the full logic of the log hook.
	// This method is called in parallel, and each log is called once.
	// Errors returned by this method are handled by the internal error handler (output to standard
	// error by default), and any errors returned will not prevent log writing.
	// The log summary instance received by this method will be cleaned up and recycled after the
	// log writing is completed. The log hook program should not hold the log summary instance.
	Fire(Summary) error
}

// NewHookFromFunc returns a log hook created from a given log level and function.
func NewHookFromFunc(levels []Level, hook func(Summary) error) Hook {
	return &hookFuncWrapper{levels, hook}
}

// The hookFuncWrapper type is used to wrap a given list of log levels and a function as a log hook.
type hookFuncWrapper struct {
	levels []Level
	hook   func(Summary) error
}

// Levels returns the log levels associated with the current log hook.
func (w *hookFuncWrapper) Levels() []Level {
	return w.levels
}

// Fire receives the summary of the log and performs the full logic of the log hook.
func (w *hookFuncWrapper) Fire(s Summary) error {
	return w.hook(s)
}

// HookBag interface defines a collection of multiple log hooks.
type HookBag interface {
	Hook

	// Add adds the given log hook to the current hook bag.
	Add(Hook)
}

// NewHookBag returns a built-in HookBag instance.
func NewHookBag() HookBag {
	return &hookBag{hooks: make(map[Level][]Hook)}
}

// The hookBag type is a built-in implementation of the HookBag interface.
type hookBag struct {
	hooks map[Level][]Hook
}

// Add adds the given log hook to the current hook bag.
func (o *hookBag) Add(hook Hook) {
	for _, level := range hook.Levels() {
		if level.IsValid() {
			o.hooks[level] = append(o.hooks[level], hook)
		}
	}
}

// Levels returns the log levels associated with the current log hook.
func (o *hookBag) Levels() []Level {
	r := make([]Level, 0, len(o.hooks))
	for level := range o.hooks {
		r = append(r, level)
	}
	return r
}

// Fire receives the summary of the log and performs the full logic of the log hook.
func (o *hookBag) Fire(s Summary) error {
	if len(o.hooks) == 0 || len(o.hooks[s.Level()]) == 0 {
		return nil
	}
	hs := o.hooks[s.Level()]
	for i, j := 0, len(hs); i < j; i++ {
		if err := hs[i].Fire(s); err != nil {
			return err
		}
	}
	return nil
}
