package logging

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

type FakeLogger struct {
	C chan *LogItem

	log_level atomic.Pointer[LogLevel]

	out  io.Writer
	lock sync.Mutex

	prefix atomic.Pointer[string] // prefix on each line to identify the logger (but see Lmsgprefix)
	flag   atomic.Int32           // properties

}

func (l *FakeLogger) Fatal(v ...any) {
	tmp := &LogItem{
		Time:    time.Now(),
		Level:   LogLevelFatal,
		Message: *l.prefix.Load() + fmt.Sprint(v...),
	}
	l.C <- tmp

}
func (l *FakeLogger) Fatalf(format string, v ...any) {
	tmp := &LogItem{
		Time:    time.Now(),
		Level:   LogLevelFatal,
		Message: *l.prefix.Load() + fmt.Sprint(v...),
	}
	l.C <- tmp
}
func (l *FakeLogger) Fatalln(v ...any) {
	tmp := &LogItem{
		Time:    time.Now(),
		Level:   LogLevelFatal,
		Message: *l.prefix.Load() + fmt.Sprint(v...),
	}
	l.C <- tmp
}
func (l *FakeLogger) Flags() int {
	return int(l.flag.Load())
}
func (l *FakeLogger) Output(calldepth int, s string) error {
	return nil
}
func (l *FakeLogger) Panic(v ...any) {
	tmp := &LogItem{
		Time:    time.Now(),
		Level:   LogLevelError02,
		Message: *l.prefix.Load() + fmt.Sprint(v...),
	}
	l.C <- tmp
	panic(tmp)
}
func (l *FakeLogger) Panicf(format string, v ...any) {
	tmp := &LogItem{
		Time:    time.Now(),
		Level:   LogLevelError02,
		Message: *l.prefix.Load() + fmt.Sprintf(format, v...),
	}
	l.C <- tmp
	panic(tmp)
}
func (l *FakeLogger) Panicln(v ...any) {
	tmp := &LogItem{
		Time:    time.Now(),
		Level:   LogLevelError02,
		Message: *l.prefix.Load() + fmt.Sprint(v...),
	}
	l.C <- tmp
	panic(tmp)
}
func (l *FakeLogger) Prefix() string {
	return *l.prefix.Load()
}
func (l *FakeLogger) Print(v ...any) {
	tmp := &LogItem{
		Time:    time.Now(),
		Level:   LogLevelInfo01,
		Message: *l.prefix.Load() + fmt.Sprint(v...),
	}
	l.C <- tmp
}
func (l *FakeLogger) Printf(format string, v ...any) {
	tmp := &LogItem{
		Time:    time.Now(),
		Level:   LogLevelInfo01,
		Message: *l.prefix.Load() + fmt.Sprintf(format, v...),
	}
	l.C <- tmp
}
func (l *FakeLogger) Println(v ...any) {
	tmp := &LogItem{
		Time:    time.Now(),
		Level:   LogLevelInfo01,
		Message: *l.prefix.Load() + fmt.Sprint(v...),
	}
	l.C <- tmp
}
func (l *FakeLogger) SetFlags(flag int) {
	l.flag.Store(int32(flag))
}
func (l *FakeLogger) SetOutput(w io.Writer) {
	l.lock.Lock()
	l.out = w
	l.lock.Unlock()
}
func (l *FakeLogger) SetPrefix(prefix string) {
	l.prefix.Store(&prefix)
}
func (l *FakeLogger) Writer() io.Writer {
	return &FakeWriter{
		f: l,
	}
}

type FakeWriter struct {
	f *FakeLogger
}

func (fw *FakeWriter) Write(p []byte) (n int, err error) {

	fw.f.lock.Lock()
	n, err = fw.f.out.Write(p)
	fw.f.lock.Unlock()

	tmp := &LogItem{
		Time:    time.Now(),
		Message: "writing: " + string(p),
		Level:   LogLevelWrite,
	}

	if err != nil {
		tmp.CallStack = []error{err}
	}
	fw.f.C <- tmp

	return
}
