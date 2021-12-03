package logger

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestNewConsoleFormatter(t *testing.T) {
	if f := NewConsoleFormatter(); f == nil {
		t.Fatal("NewConsoleFormatter(): return nil.")
	}
}

func TestConsoleFormatter_Format(t *testing.T) {
	l := New("CONSOLE")
	l.SetFormatter(NewConsoleFormatter())
	l.SetDefaultTimeFormat("TIME")
	l.SetExitFunc(nil)
	l.SetPanicFunc(nil)
	buf := new(bytes.Buffer)
	l.SetOutput(buf)

	l.WithField("foo", 1).Trace("test")
	l.WithField("foo", 1).Debug("test")
	l.WithField("foo", 1).Info("test")
	l.WithField("foo", 1).Warn("test")
	l.WithField("foo", 1).Error("test")
	l.WithField("foo", 1).Fatal("test")
	l.WithField("foo", 1).Panic("test")
	fmt.Println(buf.String())
	lines := []string{
		"CONSOLE [TIME][\u001B[36mTAC\u001B[0m] test foo=1",
		"CONSOLE [TIME][\u001B[96mDBG\u001B[0m] test foo=1",
		"CONSOLE [TIME][\u001B[92mINF\u001B[0m] test foo=1",
		"CONSOLE [TIME][\u001B[93mWAN\u001B[0m] test foo=1",
		"CONSOLE [TIME][\u001B[95mERR\u001B[0m] test foo=1",
		"CONSOLE [TIME][\u001B[31mFAT\u001B[0m] test foo=1",
		"CONSOLE [TIME][\u001B[91mPNC\u001B[0m] test foo=1",
	}
	if want := strings.Join(lines, "\n") + "\n"; want != buf.String() {
		t.Fatalf("NewConsoleFormatter().Format(): %s", strconv.Quote(buf.String()))
	}
}

func TestConsoleFormatter_Format_WithoutName(t *testing.T) {
	l := New("")
	l.SetFormatter(NewConsoleFormatter())
	l.SetDefaultTimeFormat("TIME")
	buf := new(bytes.Buffer)
	l.SetOutput(buf)

	l.WithField("foo", 1).Trace("test")
	fmt.Println(buf.String())
	want := "[TIME][\u001B[36mTAC\u001B[0m] test foo=1\n"
	if want != buf.String() {
		t.Fatalf("NewConsoleFormatter().Format(): %s", strconv.Quote(buf.String()))
	}
}
