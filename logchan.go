
package logchan

import (
	"log"
	"fmt"
	"strings"
)

type Level uint64

type Channel struct {
	Level Level
	Key byte
	Desc string
}

type Channels []Channel

type Logger struct {
	level Level
	channels Channels
	chanmap map[byte]Channel
	bitmap map[Level]Channel
}

func NewLogger(ch Channels, def Level) *Logger {
	ret := new(Logger)
	ret.level = def
	ret.chanmap = make(map[byte]Channel)
	ret.bitmap= make(map[Level]Channel)
	for _, c := range ch {
		ret.chanmap[c.Key] = c
		ret.bitmap[c.Level] = c
	}
	ret.channels = ch
	return ret
}


func (logger *Logger) AtLevel (l Level) bool {
	return (l & logger.level) != 0
}

func (logger *Logger) LevelToPrefix (l Level) string {
	s := logger.LevelToString(l)
	if len(s) > 0 {
		s = fmt.Sprintf("[%s] ", s)
	}
	return s
}

func (logger *Logger) LevelToString(l Level) string {
	descs := make([]string,0)
	
	
	for i := 0; l != 0 && i < len(logger.channels); i++ {
		c := logger.channels[i]
		if (c.Level & l) != 0 {
			descs = append (descs, c.Desc)
			l = l&(^c.Level)
		}
	}

	return strings.Join(descs, ",")
}

func (logger *Logger) SetChannels (s string) (e error, newdesc string) {

	var newlev Level = 0


	for _, c := range []byte(s) {
		if ch, found := logger.chanmap[c]; found {
			newlev |= ch.Level
		} else {
			e = fmt.Errorf("bad logger channel found: '%c'\n", c)
			break
		}
	}

	if e == nil {
		newdesc = logger.LevelToString(newlev)
		logger.level = newlev
	}
	return
}


func (logger *Logger) Printf(l Level, fmt string, v ...interface{}) {
	if logger.AtLevel (l) {
		s := logger.LevelToPrefix (l)
		if len(s) > 0 {
			fmt = s + fmt
		}
		log.Printf (fmt, v...)
	}
}

func (logger *Logger) Print(l Level, v ...interface{}) {
	if logger.AtLevel (l) {
		s := logger.LevelToPrefix (l)
		if len(s) > 0 {
			log.Print (s)
		}
		log.Print(v)
	}
}

func (logger *Logger) Println(l Level, v ...interface{}) {
	if logger.AtLevel (l) {
		s := logger.LevelToPrefix (l)
		if len(s) > 0 {
			log.Print (s)
		}
		log.Println (v...)
	}
}

