package logger

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
)

func TestConsoleFormatter_Format_WithoutCaller(t *testing.T) {
	l := New("")
	l.SetFormatter(NewConsoleFormatter())
	l.SetDefaultTimeFormat("TIME")
	buf := new(bytes.Buffer)
	l.SetOutput(buf)

	l.WithField("foo", 1).WithCaller().Trace("test") // LINE 17
	fmt.Println(buf.String())
	want := "[TIME][\u001B[36mTAC\u001B[0m] test console_formatter_caller_test.go:17 foo=1\n"
	if want != buf.String() {
		t.Fatalf("NewConsoleFormatter().Format(): %s", strconv.Quote(buf.String()))
	}
}
