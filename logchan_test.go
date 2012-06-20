package logchan

import (
	"fmt"
	"log"
	"strings"
	"testing"
	)

type ArrayWriter struct {
	lines [][]byte
}

func NewArrayWriter() *ArrayWriter {
	return &ArrayWriter{make([][]byte, 0)}
}

func (w *ArrayWriter) Write(line []byte) (n int, err error) {
	w.lines = append(w.lines, line)
	n = len(line)
	return
}

func (w *ArrayWriter) String() string {
	strs := make([]string, 0)
	for _, l := range w.lines {
		strs = append(strs, string(l))
	}
	return strings.Join(strs, "\n")
}

func (w *ArrayWriter) Len() int {
	return len(w.lines)
}

func TestStdLog(t *testing.T) {
	w := NewArrayWriter()
	log.SetOutput(w)
	SetChannels("A")
	if w.Len() != 0 {
		t.Errorf("Newly created ArrayWriter has non-zero length: %v\n", w.Len())
	}

	Printf(LOG_FATAL, "foo bar")

	if w.Len() != 1 {
		t.Errorf("Failed to log 'foo bar'")
	}
		
	fmt.Printf("%s", w)
}