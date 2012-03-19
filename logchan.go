
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
	channels Channels
}

func (logger *Logger) AtLevel (l Level) bool {
	return (l & logger.level) == l
}

func (logger *Logger) SetLogger(s string) (e error, newdesc string) {

	tab := make(map[byte]Channel)

	for _, c := range logger.channels {
		tab[c.key] = c
	}

	var newlev Level = 0


	descs := make([]string,len(s))

	for _, c := range []byte(s) {
		if ch, found := tab[c]; found {
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

