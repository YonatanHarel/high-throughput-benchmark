package logging

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

type Config struct {
	Level 			string
	Format 			string
	DisableReqLogs	bool
}

type Logger struct {
	level 			Level
	json 			bool
	disableReqLogs	bool
	logger			*log.Logger
}

var std = &Logger{
	level: Info,
	logger: log.New(os.Stdout, "", log.LstdFlags),
}

func Init(cfg Config) {
	lvl := parseLevel(cfg.Level)
	json := strings.ToLower(cfg.Format) == "json"

	std = &Logger{
		level: 			lvl,
		json: 			json,
		disableReqLogs:	cfg.DisaleReqLogs,
		logger: 		log.New(os.Stdout, "", log.LstdFlags),
	}
}


func parseLevel(s string) Level {
	switch strings.ToLower(s) {
		case "debug":
			return Debug
		case "info":
			return Info
		case "warn":
			return Warn
		case "error":
			return Error
		default:
			return Info
	}	
}

func (l *Logger) log(level Level, tag string, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	msg := fmt.Sprintf(format, args...)

	if l.json {
		l.logger.Printf(`{"level":"%s", "message":"%s"}`,  tag, msg)
	} else {
		l.logger.Printf("[%s] %s", tag, msg)
	}
}

func Debugf(format string, args... interface{}) {
	std.log(Debug, "DEBUG", format, args...)
}

func Infof(format string, args... interface{}) {
	std.log(Info, "INFO", format, args...)
}

func Warnf(format string, args... interface{}) {
	std.log(Warn, "WARN", format, args...)
}	

func Errorf(format string, args... interface{}) {
	std.log(Error, "ERROR", format, args...)
}

func requestf(format string, args... interface{}) {
	if std.disableReqLogs {
		return
	}
	std.log(Info, "REQUEST", format, args...)
}

func DisableReqLogs() bool {
	return std.disableReqLogs
}