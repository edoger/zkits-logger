# ZKits Logger Library #

[English](README.md) | 简体中文

[![ZKits](https://img.shields.io/badge/ZKits-Library-f3c)](https://github.com/edoger/zkits-logger)
[![Build Status](https://travis-ci.org/edoger/zkits-logger.svg?branch=master)](https://travis-ci.org/edoger/zkits-logger)
[![Build status](https://ci.appveyor.com/api/projects/status/xpbbppv3aui8n3fb/branch/master?svg=true)](https://ci.appveyor.com/project/edoger56924/zkits-logger/branch/master)
[![Coverage Status](https://coveralls.io/repos/github/edoger/zkits-logger/badge.svg?branch=master)](https://coveralls.io/github/edoger/zkits-logger?branch=master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/cb497bf703f44950afb43b51b3a0e581)](https://www.codacy.com/manual/edoger/zkits-logger?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=edoger/zkits-logger&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/edoger/zkits-logger)](https://goreportcard.com/report/github.com/edoger/zkits-logger)
[![Golang Version](https://img.shields.io/badge/golang-1.13+-orange)](https://github.com/edoger/zkits-logger)

## 简介 ##

这个库是 ZKits 项目的一部分，我们在这里提供了一个完整的零依赖的JSON日志库，与日志标准库完全兼容。

- 灵活可控的日志 Caller 报告，支持按日志级别报告。
- 支持 7 种标准的日志级别。
- 与 Golang 标准库完全兼容。
- 链式调用，支持为每条日志添加扩展字段，更方便排查应用程序的问题。
- 灵活的日志钩子支持。
- 高度可定制的日志格式，可自由配置日志格式化器。

## 安装 ##

```sh
go get -u -v github.com/edoger/zkits-logger
```

## 使用指南 ##

```go
package main

import (
    "github.com/edoger/zkits-logger"
)

func main() {
    // 创建一个指定名称的日志记录器。
    log := logger.New("test")

    // {"level":"info","message":"Something happened.","name":"test","time":"2020-02-20T20:20:20+08:00"}
    log.Info("Something happened.")

    // {"fields":{"num":1},"level":"info","message":"Something happened.","name":"test","time":"2020-02-20T20:20:20+08:00"}
    log.WithField("num", 1).Info("Something happened.")
}
```

## 许可证 ##

[Apache-2.0](http://www.apache.org/licenses/LICENSE-2.0)
