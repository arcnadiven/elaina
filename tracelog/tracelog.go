package tracelog

import (
	"fmt"
	"strings"
)

const (
	prefix_format = " %s(%s) "
)

type TraceLogger interface {
	WithValue(key, value string)
	CleanUp()

	// TODO: Add log level yourself if you need more
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
}

type TraceLoggerImpl struct {
	logger BaseLogger
	cache  []traceMessage // change map to list for sequences
	prefix string
}

type traceMessage struct {
	Key, Value string
}

func NewTraceLogger(baseLog BaseLogger) TraceLogger {
	return &TraceLoggerImpl{
		logger: baseLog,
		cache:  []traceMessage{},
		prefix: "",
	}
}

func (l *TraceLoggerImpl) WithValue(key, value string) {
	isChanged, isExist := true, false
	for idx, msg := range l.cache {
		if msg.Key == key {
			isExist = true
			if msg.Value == value {
				isChanged = false
			}
			l.cache[idx].Value = value
		}
	}
	if !isExist {
		l.cache = append(l.cache, traceMessage{Key: key, Value: value})
	}

	// update prefix
	if isChanged {
		tvList := []string{}
		for _, msg := range l.cache {
			tvList = append(tvList, fmt.Sprintf(prefix_format, msg.Key, msg.Value))
		}
		l.prefix = strings.Join(tvList, ",")
	}
}

func (l *TraceLoggerImpl) CleanUp() {
	l.cache = []traceMessage{}
	l.prefix = ""
}

func (l *TraceLoggerImpl) Infof(format string, args ...interface{}) {
	l.logger.Infof(l.prefix+format, args...)
}

func (l *TraceLoggerImpl) Infoln(args ...interface{}) {
	l.logger.Infoln(append([]interface{}{l.prefix}, args...)...)
}

func (l *TraceLoggerImpl) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(l.prefix+format, args...)
}

func (l *TraceLoggerImpl) Errorln(args ...interface{}) {
	l.logger.Errorln(append([]interface{}{l.prefix}, args...)...)
}
