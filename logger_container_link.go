// Author: Steve Zhang
// Date: 2020/9/16 4:55 下午

package log

import "github.com/sirupsen/logrus"

func (ct *LoggerContainer) Log(level logrus.Level, fields logrus.Fields) {
	logger := ct.MustGetLogger()
	defer ct.PutLogger(logger)

	entry := logger.WithFields(fields)

	switch level {
	case logrus.InfoLevel:
		entry.Info()
	case logrus.WarnLevel:
		entry.Warn()
	case logrus.ErrorLevel:
		entry.Error()
	case logrus.FatalLevel:
		entry.Fatal()
	case logrus.PanicLevel:
		entry.Panic()
	}
}

func (ct *LoggerContainer) Info(fields map[string]interface{}) {
	ct.Log(logrus.InfoLevel, fields)
}

func (ct *LoggerContainer) Warn(fields map[string]interface{}) {
	ct.Log(logrus.WarnLevel, fields)
}

func (ct *LoggerContainer) Error(fields map[string]interface{}) {
	ct.Log(logrus.ErrorLevel, fields)
}

func (ct *LoggerContainer) Debug(fields map[string]interface{}) {
	ct.Log(logrus.DebugLevel, fields)
}

func (ct *LoggerContainer) Fatal(fields map[string]interface{}) {
	ct.Log(logrus.FatalLevel, fields)
}

func (ct *LoggerContainer) Panic(fields map[string]interface{}) {
	ct.Log(logrus.PanicLevel, fields)
}
