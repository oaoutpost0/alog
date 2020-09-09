// Written by Gon Yi
// This is a modified version of Go's standard logger -- trimmed down and added category logging.

package alog

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type FormatFlag int

const (
	// FLAG FORMAT
	F_ALL  FormatFlag = ^FormatFlag(0)
	F_NONE FormatFlag = 0
	F_DATE FormatFlag = 1 << iota
	F_TIME
	F_MICROSECONDS
	F_UTC
	F_PREFIX

	// FLAG CATEGORY
	ALL     uint64 = ^uint64(0)
	C_ALL          = ALL
	C_NONE  uint64 = 0
	C_DEBUG uint64 = 1 << iota
	C_INFO
	C_WARN
	C_ERROR
	C_FATAL
	C_SYSTEM
	C_IO
	C_NET
	C_SERVICE
	C_REQ // (service) SEND the request
	C_RES // (service) RECEIVEd request and responsed
	C_TIMED
	C_API
	C_BACK
	C_FRONT
	C_GLOBAL
)

type Alog struct {
	mu         sync.Mutex
	prefix     string
	flagFormat FormatFlag
	flagCat    uint64
	out        io.Writer
	buf        []byte

	bAltOut   bool
	altOutput func(flagCat uint64, s string) error
}

type Opt func(*Alog)

func New(opts ...Opt) *Alog {
	l := &Alog{out: os.Stderr, prefix: "", flagFormat: F_TIME | F_PREFIX, flagCat: C_ALL,
		altOutput: func(uint64, string) error { return nil },
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *Alog) SetAltOutput(fnAltOut func(uint64, string) error) {
	if fnAltOut == nil {
		l.bAltOut = false
		l.altOutput = func(uint64, string) error { return nil }
	} else {
		l.bAltOut = true
		l.altOutput = fnAltOut
	}
}

func (l *Alog) SetOutput(w io.Writer) {
	l.mu.Lock()
	l.out = w
	l.mu.Unlock()
}

func itoa(buf *[]byte, i int, wid int) {
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (l *Alog) formatHeader(buf *[]byte, t time.Time) {
	if l.flagFormat&(F_DATE|F_TIME|F_MICROSECONDS) != 0 {
		if l.flagFormat&F_UTC != 0 {
			t = t.UTC()
		}
		if l.flagFormat&F_DATE != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flagFormat&(F_TIME|F_MICROSECONDS) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flagFormat&F_MICROSECONDS != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.flagFormat&F_PREFIX != 0 {
		*buf = append(*buf, l.prefix...)
	}
}

func (l *Alog) Output(flagCat uint64, s string) error {
	if l.flagCat&flagCat != 0 { // if category matches
		if l.bAltOut {
			return l.altOutput(flagCat, s)
		}

		now := time.Now() // get this early.
		l.mu.Lock()
		l.buf = l.buf[:0]
		l.formatHeader(&l.buf, now)
		l.buf = append(l.buf, s...)
		if len(s) == 0 || s[len(s)-1] != '\n' {
			l.buf = append(l.buf, '\n')
		}
		_, err := l.out.Write(l.buf)
		l.mu.Unlock()
		return err
	}
	return nil
}

// Outputb is for exported writer as interface uses func([]byges)(int,error)
func (l *Alog) Outputb(flagCat uint64, b []byte) error {
	if l.flagCat&flagCat != 0 { // if category matches
		now := time.Now() // get this early.
		l.mu.Lock()
		l.buf = l.buf[:0]
		l.formatHeader(&l.buf, now)
		l.buf = append(l.buf, b...)
		if len(b) == 0 || b[len(b)-1] != '\n' {
			l.buf = append(l.buf, '\n')
		}
		_, err := l.out.Write(l.buf)
		l.mu.Unlock()
		return err
	}
	return nil
}

// ==================================================================== PRINT
func (l *Alog) Print(flagCat uint64, s string)           { l.Output(flagCat, s) }
func (l *Alog) Println(flagCat uint64, v ...interface{}) { l.Output(flagCat, fmt.Sprintln(v...)) }
func (l *Alog) Printf(flagCat uint64, format string, v ...interface{}) {
	l.Output(flagCat, fmt.Sprintf(format, v...))
}

// ==================================================================== IFACE
func (l *Alog) Debug(s string) {
	l.Output(C_DEBUG, s)
}
func (l *Alog) Debugf(f string, v ...interface{}) {
	l.Output(C_DEBUG, fmt.Sprintf(f, v...))
}
func (l *Alog) Info(s string) {
	l.Output(C_INFO, s)
}
func (l *Alog) Infof(f string, v ...interface{}) {
	l.Output(C_INFO, fmt.Sprintf(f, v...))
}
func (l *Alog) Warn(s string) {
	l.Output(C_WARN, s)
}
func (l *Alog) Warnf(f string, v ...interface{}) {
	l.Output(C_WARN, fmt.Sprintf(f, v...))
}
func (l *Alog) Error(s string) {
	l.Output(C_ERROR, s)
}
func (l *Alog) Errorf(f string, v ...interface{}) {
	l.Output(C_ERROR, fmt.Sprintf(f, v...))
}

// ==================================================================== FATAL
func (l *Alog) Fatal(flagCat uint64, s string) {
	l.Output(flagCat, s)
	os.Exit(1)
}
func (l *Alog) Fatalln(flagCat uint64, v ...interface{}) {
	l.Output(flagCat, fmt.Sprintln(v...))
	os.Exit(1)
}
func (l *Alog) Fatalf(flagCat uint64, format string, v ...interface{}) {
	l.Output(flagCat, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// ==================================================================== PANIC
func (l *Alog) Panic(flagCat uint64, s string) {
	l.Output(flagCat, s)
	panic(s)
}
func (l *Alog) Panicln(flagCat uint64, v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.Output(flagCat, s)
	panic(s)
}
func (l *Alog) Panicf(flagCat uint64, format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Output(flagCat, s)
	panic(s)
}

// ==================================================================== FLAG
func (l *Alog) SetFormat(flagFmt ...FormatFlag) {
	var flag FormatFlag = 0
	switch len(flagFmt) {
	case 0:
		flag = F_NONE
	case 1:
		flag = flagFmt[0]
	default:
		for _, v := range flagFmt {
			flag = flag | v
		}
	}
	l.mu.Lock()
	l.flagFormat = flag
	l.mu.Unlock()
}
func (l *Alog) SetFilter(category ...uint64) {
	var flag uint64 = 0

	switch len(category) {
	case 0:
		flag = C_NONE
	case 1:
		flag = category[0]
	default:
		for _, v := range category {
			flag = flag | v
		}
	}

	l.mu.Lock()
	l.flagCat = flag
	l.mu.Unlock()
}

// ==================================================================== PREFIX
func (l *Alog) Prefix() string {
	l.mu.Lock()
	l.mu.Unlock()
	return l.prefix
}
func (l *Alog) SetPrefix(prefix string) {
	l.mu.Lock()
	l.prefix = prefix
	l.mu.Unlock()
}

// ==================================================================== WRITER
func (l *Alog) Writer() io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out
}
func (l *Alog) NewWriter(flag uint64, prefix string) io.Writer {
	if prefix != "" {
		return &LogWriterPrefix{
			alog:       l,
			flagCat:    flag,
			prefix:     []byte(prefix),
			prefixSize: len(prefix),
			usePrefix:  true,
		}
	}
	return &LogWriter{
		alog:    l,
		flagCat: flag,
	}
}

// ==================================================================== PRINT FUNC
func (l *Alog) NewPrint(flagCat uint64) func(string) {
	return func(s string) {
		l.Output(flagCat, s)
	}
}
func (l *Alog) NewPrintln(flagCat uint64) func(...interface{}) {
	return func(v ...interface{}) {
		l.Output(flagCat, fmt.Sprint(v...))
	}
}

func (l *Alog) NewPrintf(flagCat uint64) func(string, ...interface{}) {
	return func(s string, v ...interface{}) {
		l.Output(flagCat, fmt.Sprintf(s, v...))
	}
}

// ==================================================================== MINI DISCARD
type devNull int

var Discard io.Writer = devNull(0)

func (devNull) Write(p []byte) (int, error) {
	return 0, nil
}

// ==================================================================== WRITER
type LogWriter struct {
	alog    *Alog
	flagCat uint64
}

func (dw *LogWriter) Write(p []byte) (n int, err error) {
	dw.alog.Outputb(dw.flagCat, p)
	return len(p), nil
}

// ==================================================================== WRITER: PREFIX
type LogWriterPrefix struct {
	alog       *Alog
	flagCat    uint64
	prefix     []byte
	prefixSize int
	usePrefix  bool
}

func (dw *LogWriterPrefix) Write(p []byte) (n int, err error) {
	if dw.usePrefix {
		dw.alog.Outputb(dw.flagCat, append(dw.prefix, p...))
		return len(p) + dw.prefixSize, nil
	}
	dw.alog.Outputb(dw.flagCat, p)
	return len(p), nil
}

// ==================================================================== GLOBAL
var std = New(func(l *Alog) {
	l.SetFilter(C_ALL)
	l.SetFormat()
	l.SetOutput(os.Stderr)
})

func NewPrint(flagCat uint64) func(string)                  { return std.NewPrint(flagCat) }
func NewPrintln(flagCat uint64) func(...interface{})        { return std.NewPrintln(flagCat) }
func NewPrintf(flagCat uint64) func(string, ...interface{}) { return std.NewPrintf(flagCat) }
func NewWriter(flag uint64, prefix string) io.Writer        { return std.NewWriter(flag, prefix) }

func SetOutput(w io.Writer)           { std.SetOutput(w) }
func SetPrefix(prefix string)         { std.SetPrefix(prefix) }
func SetFormat(flagFmt ...FormatFlag) { std.SetFormat(flagFmt...) }
func SetFilter(flagCat uint64)        { std.SetFilter(flagCat) }

func Print(s string)                         { std.Print(C_GLOBAL, s) }
func Println(v ...interface{})               { std.Println(C_GLOBAL, v...) }
func Printf(format string, v ...interface{}) { std.Printf(C_GLOBAL, format, v...) }

func Fatal(s string)                         { std.Fatal(C_GLOBAL, s) }
func Fatalln(v ...interface{})               { std.Fatalln(C_GLOBAL, v...) }
func Fatalf(format string, v ...interface{}) { std.Fatalf(C_GLOBAL, format, v...) }

func Panic(s string)                         { std.Panic(C_GLOBAL, s) }
func Panicln(v ...interface{})               { std.Panicln(C_GLOBAL, v...) }
func Panicf(format string, v ...interface{}) { std.Panicf(C_GLOBAL, format, v...) }

func Debug(s string)                    { std.Debug(s) }
func Debugf(f string, v ...interface{}) { std.Debugf(f, v...) }
func Info(s string)                     { std.Info(s) }
func Infof(f string, v ...interface{})  { std.Infof(f, v...) }
func Warn(s string)                     { std.Warn(s) }
func Warnf(f string, v ...interface{})  { std.Warnf(f, v...) }
func Error(s string)                    { std.Error(s) }
func Errorf(f string, v ...interface{}) { std.Errorf(f, v...) }
