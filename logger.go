package log

import (
	"fmt"
	"io"
	"time"

	"github.com/rifflock/lfshook"

	rotate "github.com/lestrrat/go-file-rotatelogs"

	"github.com/sirupsen/logrus"
)

const timeLayout = "2006-01-02 15:04:05"

type LoggerConf struct {
	Level        string // 日志级别
	ReportCaller bool   // 是否打印调用者

	InfLinkName    string // 通知日志文件链接名称
	InfMaxAgeHours int    // 通知日志最大保存小时数
	InfRotateHours int    // 通知日志轮转小时数

	ErrLinkName    string // 错误日志链接名称
	ErrMaxAgeHours int    // 错误日志文件最大保存小时数
	ErrRotateHours int    // 错误日志文件轮转小时数
}

type Logger struct {
	*logrus.Logger

	infout, errout io.Closer
}

func NewLogger(cf *LoggerConf) (logger *Logger, err error) {
	logger = &Logger{
		Logger: logrus.New(),
	}

	// 设置标准输出格式
	logger.SetFormatter()

	// 设置是否开启调用者打印
	logger.SetReportCaller(cf)

	// 设置日志级别
	if err = logger.SetLevel(cf); err != nil {
		err = fmt.Errorf("set level: %w", err)
		return
	}

	// 设置勾子
	if err = logger.SetHook(cf); err != nil {
		err = fmt.Errorf("set hook: %w", err)
		return
	}

	return
}

func (logger *Logger) Close() (err error) {
	var closeInfErr, closeErrErr error

	if logger.infout != nil {
		closeInfErr = logger.infout.Close()
	}

	if logger.errout != nil {
		closeErrErr = logger.errout.Close()
	}

	if closeInfErr != nil || closeErrErr != nil {
		err = fmt.Errorf("close infout: %v, close errout: %v", closeInfErr, closeErrErr)
		return
	}

	return
}

func (logger *Logger) SetFormatter() {
	logger.Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: timeLayout,
	})
}

func (logger *Logger) SetReportCaller(cf *LoggerConf) {
	logger.Logger.SetReportCaller(cf.ReportCaller)
}

func (logger *Logger) SetLevel(cf *LoggerConf) (err error) {
	level, err := logrus.ParseLevel(cf.Level)
	if err != nil {
		err = fmt.Errorf("parse level: %w", err)
		return
	}
	logger.Logger.SetLevel(level)
	return
}

func (logger *Logger) SetHook(cf *LoggerConf) (err error) {
	infout, err := rotate.New(
		cf.InfLinkName+".%Y-%m-%d-%H", rotate.WithLinkName(cf.InfLinkName),
		rotate.WithMaxAge(time.Duration(cf.InfMaxAgeHours)*time.Hour),
		rotate.WithRotationTime(time.Duration(cf.InfRotateHours)*time.Hour),
	)
	if err != nil {
		err = fmt.Errorf("new info rotate: %w", err)
		return
	}

	errout, err := rotate.New(
		cf.ErrLinkName+".%Y-%m-%d-%H", rotate.WithLinkName(cf.ErrLinkName),
		rotate.WithMaxAge(time.Duration(cf.ErrMaxAgeHours)*time.Hour),
		rotate.WithRotationTime(time.Duration(cf.ErrRotateHours)*time.Hour),
	)
	if err != nil {
		err = fmt.Errorf("new error rotate: %w", err)
		return
	}

	hook := lfshook.NewHook(
		lfshook.WriterMap{
			logrus.DebugLevel: infout,
			logrus.InfoLevel:  infout,
			logrus.WarnLevel:  errout,
			logrus.ErrorLevel: errout,
			logrus.FatalLevel: errout,
			logrus.PanicLevel: errout,
		},
		&logrus.JSONFormatter{
			TimestampFormat: timeLayout,
		},
	)

	hooks := make(logrus.LevelHooks)
	hooks.Add(hook)

	logger.Logger.ReplaceHooks(hooks)

	logger.infout = infout
	logger.errout = errout

	return
}
