
package logchan

import (
	"log"
	"fmt"
	"strings"
)

type Level uint64

type Channel struct {
	level Level
	key byte
	desc string
}

type Channels []Channel

type Logger struct {
	level Level
	chanmap map[byte]Channel
}

func NewLogger(ch Channels, def Level) *Logger {
	var ret *Logger
	ret.level = def
	ret.chanmap = make(map[byte]Channel)
	for _, c := range ch {
		ret.chanmap[c.key] = c
	}
	return ret
}


func (logger *Logger) AtLevel (l Level) bool {
	return (l & logger.level) == l
}

func (logger *Logger) SetChannels (s string) (e error, newdesc string) {

	var newlev Level = 0


	descs := make([]string,len(s))

	for _, c := range []byte(s) {
		if ch, found := logger.chanmap[c]; found {
			newlev |= ch.level
			descs = append (descs, ch.desc)
		} else {
			e = fmt.Errorf("bad logger channel found: '%c'\n", c)
			break
		}
	}

	if e == nil {
		newdesc = strings.Join(descs, ",")
		logger.level = newlev
	}
	return
}


func (logger *Logger) Printf(l Level, fmt string, v ...interface{}) {
	if logger.AtLevel (l) {
		log.Printf (fmt, v...)
	}
}

func (logger *Logger) Print(l Level, v ...interface{}) {
	if logger.AtLevel (l) {
		log.Print (v...)
	}
}

func (logger *Logger) Println(l Level, v ...interface{}) {
	if logger.AtLevel (l) {
		log.Println (v...)
	}
}

