package alog_test

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/orangenumber/alog"
)

func TestALog(t *testing.T) {
	buf := bytes.Buffer{}

	alog.SetOutput(&buf)
	alog.SetFormat(0)

	// not sending any
	buf.Reset()
	alog.Print("test1")
	if r := buf.String(); r != "test1\n" {
		t.Errorf("test1: [%s]", r)
	}

	// test2\n
	buf.Reset()
	alog.Printf("test%d", 2)
	if r := buf.String(); r != "test2\n" {
		t.Errorf("test2: [%s]", r)
	}
}

func TestNewALog(t *testing.T) {
	buf := bytes.Buffer{}

	l := alog.New()
	l.SetOutput(&buf)
	l.SetFormat(0)

	// not sending any
	buf.Reset()
	l.Print(0, "test1")
	if r := buf.String(); r != "" {
		t.Errorf("test1: [%s]", r)
	}

	// test2\n
	buf.Reset()
	l.Print(alog.C_INFO, "test2")
	if r := buf.String(); r != "test2\n" {
		t.Errorf("test2: [%s]", r)
	}

	// test3\n
	buf.Reset()
	l.Print(alog.C_INFO|alog.C_NET, "test3")
	if r := buf.String(); r != "test3\n" {
		t.Errorf("test3: [%s]", r)
	}
}

func TestALogNewPrint(t *testing.T) {
	buf := bytes.Buffer{}
	alog.SetOutput(&buf)
	alog.SetFormat(0)

	pIO := alog.NewPrint(alog.C_IO)
	pNET := alog.NewPrint(alog.C_NET)

	buf.Reset()
	pIO("iotest")
	if r := buf.String(); r != "iotest\n" {
		t.Errorf("test1 io: [%s]", r)
	}

	buf.Reset()
	pNET("nettest")
	if r := buf.String(); r != "nettest\n" {
		t.Errorf("test1 net: [%s]", r)
	}

	alog.SetFilter(alog.C_INFO | alog.C_IO)
	buf.Reset()
	pIO("iotest")
	if r := buf.String(); r != "iotest\n" {
		t.Errorf("test2 io: [%s]", r)
	}

	buf.Reset()
	pNET("nettest")
	if r := buf.String(); r != "" {
		t.Errorf("test2 net: [%s]", r)
	}
}

func TestNewALogNewWriter(t *testing.T) {
	buf := bytes.Buffer{}
	l := alog.New()
	l.SetOutput(&buf)
	l.SetFormat(0)

	io := l.NewWriter(alog.C_IO, "[IO] ")
	err := l.NewWriter(alog.C_ERROR, "[ERROR] ")

	lio := log.New(io, "", 0)
	lerr := log.New(err, "", 0)

	buf.Reset()
	lio.Print("testA")
	if r := buf.String(); r != "[IO] testA\n" {
		t.Errorf("testA: [%s]", r)
	}

	buf.Reset()
	lerr.Print("testB")
	if r := buf.String(); r != "[ERROR] testB\n" {
		t.Errorf("testA: [%s]", r)
	}

}

var Result string

func outsideData(b *testing.B) struct{ ID string } {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	return struct{ ID string }{strconv.Itoa(r1.Intn(1000000000))}
}

func Benchmark_NewStdLog_Output(b *testing.B) {
	data := outsideData(b) // GET DATA FROM OUTSIDE..
	b.ResetTimer()
	b.ReportAllocs()

	l := log.New(alog.Discard, "", 0)
	for n := 0; n < b.N; n++ {
		l.Output(0, data.ID)
		// l.Print(data.ID)
	}
}
func Benchmark_NewALog_Output(b *testing.B) {
	data := outsideData(b) // GET DATA FROM OUTSIDE..
	b.ResetTimer()
	b.ReportAllocs()

	l := alog.New()
	l.SetFormat(0)
	l.SetFilter(alog.C_ALL)
	l.SetOutput(alog.Discard)
	for n := 0; n < b.N; n++ {
		l.Output(alog.C_ALL, data.ID)
	}
}
func Benchmark_NewALog_Print(b *testing.B) {
	data := outsideData(b) // GET DATA FROM OUTSIDE..
	b.ResetTimer()
	b.ReportAllocs()

	l := alog.New()
	l.SetFormat(0)
	l.SetFilter(alog.C_ALL)
	l.SetOutput(alog.Discard)
	for n := 0; n < b.N; n++ {
		l.Print(alog.C_ALL, data.ID)
	}
}

func Benchmark_NewStdLog_Print(b *testing.B) {
	data := outsideData(b) // GET DATA FROM OUTSIDE..
	b.ResetTimer()
	b.ReportAllocs()

	l := log.New(alog.Discard, "", 0)
	for n := 0; n < b.N; n++ {
		l.Print(data.ID)
	}
}

func Benchmark_NewALog_Print_Filtered(b *testing.B) {
	data := outsideData(b) // GET DATA FROM OUTSIDE..
	b.ResetTimer()
	b.ReportAllocs()

	l := alog.New()
	l.SetFormat(0)
	l.SetFilter(alog.C_ERROR | alog.C_WARN)
	l.SetOutput(alog.Discard)
	for n := 0; n < b.N; n++ {
		l.Print(alog.C_IO, data.ID)
	}
}

func Benchmark_NewALog_Errorf(b *testing.B) {
	data := outsideData(b) // GET DATA FROM OUTSIDE..
	b.ResetTimer()
	b.ReportAllocs()

	l := alog.New()
	l.SetFormat(0)
	l.SetFilter(alog.C_IO)
	l.SetOutput(alog.Discard)
	for n := 0; n < b.N; n++ {
		l.Errorf("%s", data.ID)
	}
}

func Benchmark_NewStdLog_Printf(b *testing.B) {
	data := outsideData(b) // GET DATA FROM OUTSIDE..
	b.ResetTimer()
	b.ReportAllocs()

	l := log.New(alog.Discard, "", 0)
	for n := 0; n < b.N; n++ {
		l.Printf("%s", data.ID)
	}
}

func Benchmark_NewALog_NewPrint(b *testing.B) {
	data := outsideData(b) // GET DATA FROM OUTSIDE..
	b.ResetTimer()
	b.ReportAllocs()

	l := alog.New()
	l.SetFormat(0)
	l.SetFilter(alog.C_ALL)
	l.SetOutput(alog.Discard)
	errlog := l.NewPrint(alog.C_ERROR)

	for n := 0; n < b.N; n++ {
		// l.Print(dl.C_IO, data.ID)
		errlog(data.ID)
	}
}

func Benchmark_Test(b *testing.B) {
	out := []string{"test", "this", "is", "good"}
	// out2 := "testa"
	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		// l.Print(dl.C_IO, data.ID)
		fmt.Sprint(out)
		// test1(alog.Discard, out2, out2)
	}
}
