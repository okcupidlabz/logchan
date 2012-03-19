
package logchan

import (
	"log"
	"fmt"
	"strings"
)

type LogLevel uint64

type LogChannel struct {
	logLevel LogLevel
	key byte
	desc string
}

type LogChannels []LogChannel

type Logging struct {
	logLevel LogLevel
	logChannels LogChannels
}

func (logging *Logging) AtLevel (l LogLevel) bool {
	return (l & logging.logLevel) == l
}

func (logging *Logging) SetLogging(s string) (e error, newdesc string) {

	tab := make(map[byte]LogChannel)

	for _, c := range logging.logChannels {
		tab[c.key] = c
	}

	var newlev LogLevel = 0


	descs := make([]string,len(s))

	for _, c := range []byte(s) {
		if ch, found := tab[c]; found {
			newlev |= ch.logLevel
			descs = append (descs, ch.desc)
		} else {
			e = fmt.Errorf("bad logging channel found: '%c'\n", c)
			break
		}
	}

	if e == nil {
		newdesc = strings.Join(descs, ",")
		logging.logLevel = newlev
	}
	return
}


func (logging *Logging) Printf(l LogLevel, fmt string, v ...interface{}) {
	if logging.AtLevel (l) {
		log.Printf (fmt, v...)
	}
}

func (logging *Logging) Print(l LogLevel, v ...interface{}) {
	if logging.AtLevel (l) {
		log.Print (v...)
	}
}

func (logging *Logging) Println(l LogLevel, v ...interface{}) {
	if logging.AtLevel (l) {
		log.Println (v...)
	}
}

