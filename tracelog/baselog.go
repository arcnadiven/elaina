package tracelog

import (
	"github.com/arcnadiven/atalanta/xtools"
	"github.com/sirupsen/logrus"
)

type BaseLogger interface {
	// TODO: Add log level yourself if you need more
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
}

func NewBaseLogger(logFile string) BaseLogger {
	if logFile == "" {
		return logrus.StandardLogger()
	}
	return xtools.NewFileLogrus(logFile)
}
