package tracelog

import "fmt"

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
	cache  map[string]string
	prefix string
}

func NewTraceLogger(baseLog BaseLogger) TraceLogger {
	return &TraceLoggerImpl{
		logger: baseLog,
		cache:  map[string]string{},
		prefix: "",
	}
}

func (l *TraceLoggerImpl) WithValue(key, value string) {
	isChange := false
	if val, ok := l.cache[key]; ok {
		if val != value {
			isChange = true
		}
	}
	if isChange {
		prefix := ""
		for k, v := range l.cache {
			prefix += fmt.Sprintf(prefix_format, k, v)
		}
		l.prefix = prefix
	}
}

func (l *TraceLoggerImpl) CleanUp() {
	l.cache = map[string]string{}
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
