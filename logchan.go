
package m8lt

import (
	"errors"
	"log"
	"fmt"
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

func (logging *Logging) SetLogging(s string) (e error, newdesc string) {
	tab := map[byte]LogChannel

	for _, c = range logging.logChannels {
		tab[c.key] = c
	}

	var newlev LogLevel = 0


	descs := make([]string)

	for _, c = range []bytes(s) {
		if ch, found := tab[c]; found {
			newlev |= c.logLevel
			descs = append (descs, c.desc)
		} else {
			e = fmt.Errorf("bad logging channel found: '%c'\n", c)
			break
		}
	}

	if e == nil {
		newdesc = string.Join(descs, ",")
		srv.logLevel = newlev
	}
	return
}


func (logging *Logging) Printf(l LogLevel, fmt string, v ...interface{}) {
	if (logging.LogLevel & l) == l {
		log.Printf (fmt, v...)
	}
}

func (logging *Logging) Print(l LogLevel, v ...interface{}) {
	if (logging.LogLevel & l) == l {
		log.Printf (v...)
	}
}

func (logging *Logging) Println(v ...interface{}) {
	if (logging.LogLevel & l) == l {
		log.Println (v...)
	}
}

