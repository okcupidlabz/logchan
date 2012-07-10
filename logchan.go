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

const (
	LOG_NONE  Level = 0x0
	LOG_DEBUG Level = 0x800000000000000
	LOG_INFO  Level = 0x1000000000000000
	LOG_WARN  Level = 0x2000000000000000
	LOG_ERROR Level = 0x4000000000000000
	LOG_FATAL Level = 0x8000000000000000
	LOG_ALL   Level = 0xFFFFFFFFFFFFFFFF

	LOG_LEVEL_5 Level = LOG_FATAL
	LOG_LEVEL_4 Level = LOG_ERROR | LOG_LEVEL_5
	LOG_LEVEL_3 Level = LOG_WARN | LOG_LEVEL_4
	LOG_LEVEL_2 Level = LOG_INFO | LOG_LEVEL_3
	LOG_LEVEL_1 Level = LOG_DEBUG | LOG_LEVEL_2
)

const (
	CHANNEL_NONE byte  = '0'
	CHANNEL_DEBUG byte = 'D'
	CHANNEL_INFO byte  = 'I'
	CHANNEL_WARN byte  = 'W'
	CHANNEL_ERROR byte = 'E'
	CHANNEL_FATAL byte = 'F'
	CHANNEL_LEVEL_1 byte = '1'
	CHANNEL_LEVEL_2 byte = '2'
	CHANNEL_LEVEL_3 byte = '3'
	CHANNEL_LEVEL_4 byte = '4'
	CHANNEL_LEVEL_5 byte = '5'
	CHANNEL_ALL byte   = 'A'
)

var defaultChannels = []Channel {
	Channel{LOG_NONE, CHANNEL_NONE, "none"},
	Channel{LOG_DEBUG, CHANNEL_DEBUG, "debug"},
	Channel{LOG_INFO, CHANNEL_INFO, "info"},
	Channel{LOG_WARN, CHANNEL_WARN, "warn"},
	Channel{LOG_ERROR, CHANNEL_ERROR, "error"},
	Channel{LOG_FATAL, CHANNEL_FATAL, "fatal"},
	Channel{LOG_ALL, CHANNEL_ALL, "all"}}

type Logger struct {
	level Level
	channels Channels
	chanmap map[byte]Channel
	bitmap map[Level]Channel
	catchAll Channel
}

func NewLogger(ch Channels, def Level) *Logger {
	ret := new(Logger)
	ret.level = def
	ret.channels = make(Channels, 1)
	ret.chanmap = make(map[byte]Channel)
	ret.bitmap = make(map[Level]Channel)

	l := len(defaultChannels)
	
	for i, c := range defaultChannels {
		ret.chanmap[c.Key] = c
		ret.bitmap[c.Level] = c
		ret.channels = append(ret.channels, c)
		if i != l - 1 {
			ret.channels = append(ret.channels, c)	
		} else {
			ret.catchAll = c
		}
	}

	for _, c := range ch {
		ret.chanmap[c.Key] = c
		ret.bitmap[c.Level] = c
		ret.channels = append(ret.channels, c)
	}
	return ret
}

func (logger *Logger) AddChannels(ch Channels) {
	for _, c := range ch {
		logger.chanmap[c.Key] = c
		logger.bitmap[c.Level] = c
		logger.channels = append(logger.channels, c)
	}
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

	if len(descs) == 0 {
		descs = append(descs, logger.catchAll.Desc)
	}

	return strings.Join(descs, ",")
}

func (logger *Logger) SetChannelsEasy (which string, s string, setIfEmpty bool) bool {
	ret := true
	if len(s) > 0 || setIfEmpty {
		if s, e := logger.SetChannels(s); e == nil {
			log.Printf("Setting %s logging to '%s'\n", which, s);
		} else {
			log.Printf("Failed to set %s logging: %s\n", which, e);
			ret = false
		}
	}
	return ret
}

func (logger *Logger) SetChannels (s string) (newdesc string, e error) {
	var newlev Level = 0
	for _, c := range []byte(s) {
		if ch, found := logger.chanmap[c]; found {
			newlev |= ch.Level
		} else {
			e = fmt.Errorf("bad logger channel found: '%c'", c)
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

var std = NewLogger(defaultChannels, LOG_NONE)

func AddChannels(chs Channels) {
	std.AddChannels(chs)
}

func AtLevel (l Level) bool {
	return (l & std.level) != 0
}

func LevelToPrefix (l Level) string {
	return std.LevelToPrefix(l)
}

func LevelToString(l Level) string {
	return std.LevelToString(l)
}

func SetChannelsEasy (which string, s string, setIfEmpty bool) bool {
	return std.SetChannelsEasy(which, s, setIfEmpty)
}

func SetChannels (s string) (newdesc string, e error) {
	return std.SetChannels(s)
}

func Printf(l Level, fmt string, v ...interface{}) {
	std.Printf(l, fmt, v...)
}

func Print(l Level, v ...interface{}) {
	std.Print(l, v...)
}

func Println(l Level, v ...interface{}) {
	std.Println(l, v...)
}

