# ZKits Logger Library #

[![ZKits](https://img.shields.io/badge/ZKits-Library-f3c)](https://github.com/edoger/zkits-logger)
[![Build Status](https://travis-ci.org/edoger/zkits-logger.svg?branch=master)](https://travis-ci.org/edoger/zkits-logger)
[![Build status](https://ci.appveyor.com/api/projects/status/xpbbppv3aui8n3fb/branch/master?svg=true)](https://ci.appveyor.com/project/edoger56924/zkits-logger/branch/master)
[![Coverage Status](https://coveralls.io/repos/github/edoger/zkits-logger/badge.svg?branch=master)](https://coveralls.io/github/edoger/zkits-logger?branch=master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/cb497bf703f44950afb43b51b3a0e581)](https://www.codacy.com/manual/edoger/zkits-logger?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=edoger/zkits-logger&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/edoger/zkits-logger)](https://goreportcard.com/report/github.com/edoger/zkits-logger)
[![Golang Version](https://img.shields.io/badge/golang-1.13+-orange)](https://github.com/edoger/zkits-logger)

## About ##

This package is a library of ZKits project.
This is a zero-dependency standard JSON log library that supports structured JSON logs and is compatible with the standard library.

 - Flexible and controllable caller report.
 - Support 7 log levels.
 - Complete log standard library compatibility.
 - Chained call, supporting additional log context data.
 - Flexible log hook support.
 - Custom log formatter support.

## Install ##

```sh
go get -u -v github.com/edoger/zkits-logger
```

## Usage ##

```go
package main

import (
    "github.com/edoger/zkits-logger"
)

func main() {
    // Creates a logger instance with the specified name.
    log := logger.New("test")

    // {"level":"info","message":"Something happened.","name":"test","time":"2020-02-20T20:20:20+08:00"}
    log.Info("Something happened.")

    // {"fields":{"num":1},"level":"info","message":"Something happened.","name":"test","time":"2020-02-20T20:20:20+08:00"}
    log.WithField("num", 1).Info("Something happened.")
}
```

## License ##

[Apache-2.0](http://www.apache.org/licenses/LICENSE-2.0)
