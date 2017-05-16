package log

import (
	"log"
	"fmt"
	"io"
	"os"
	//	"legitlab.letv.cn/yig/yig/helper"
)

// These flags define which text to prefix to each log entry generated by the Logger.
const (
	// Bits or'ed together to control what's printed.
	// There is no control over the order they appear (the order listed
	// here) or the format they present (as described in the comments).
	// The prefix is followed by a colon only when Llongfile or Lshortfile
	// is specified.
	// For example, flags Ldate | Ltime (or LstdFlags) produce,
	//	2009/01/23 01:23:23 message
	// while flags Ldate | Ltime | Lmicroseconds | Llongfile produce,
	//	2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

type Logger struct {
	Logger	*log.Logger
	LogLevel  int
}

func New(out io.Writer, prefix string, flag int, level int) *Logger{
	var logger Logger
	logger.LogLevel = level
	logger.Logger = log.New(out, prefix, flag)
	return &logger
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(level int, format string, v ...interface{}) {
	if l.LogLevel >= level {
		l.Logger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(level int, v ...interface{}) {
	if l.LogLevel >= level {
		l.Logger.Output(2, fmt.Sprint(v...))
	}
}

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Println(level int, v ...interface{}) {
	if l.LogLevel >= level {
		l.Logger.Output(2, fmt.Sprintln(v...))
	}
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *Logger) Fatal(level int, v ...interface{}) {
	if l.LogLevel >= level {
		l.Logger.Output(2, fmt.Sprint(v...))
	}
	os.Exit(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(level int, format string, v ...interface{}) {
	if l.LogLevel >= level {
		l.Logger.Output(2, fmt.Sprintf(format, v...))
	}
	os.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *Logger) Fatalln(level int, v ...interface{}) {
	if l.LogLevel >= level {
		l.Logger.Output(2, fmt.Sprintln(v...))
	}
	os.Exit(1)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *Logger) Panic(level int, v ...interface{}) {
	s := fmt.Sprint(v...)
	if l.LogLevel >= level {
		l.Logger.Output(2, s)
	}
	panic(s)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Logger) Panicf(level int, format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	if l.LogLevel >= level {
		l.Logger.Output(2, s)
	}
	panic(s)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *Logger) Panicln(level int, v ...interface{}) {
	s := fmt.Sprintln(v...)
	if l.LogLevel >= level {
		l.Logger.Output(2, s)
	}
	panic(s)
}
