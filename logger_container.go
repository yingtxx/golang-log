package log

import (
	"errors"
	"fmt"
	"github.com/yingtxx/golang-conf"
)

type LoggerContainer struct {
	*conf.Container
}

type GetLoggerConfFunc func() (*LoggerConf, error)

var ErrGetLoggerConfFuncIsNil = errors.New("get logger conf func is nil")

func NewLoggerContainer(getLoggerConf GetLoggerConfFunc) (ct *LoggerContainer, err error) {
	if getLoggerConf == nil {
		err = ErrGetLoggerConfFuncIsNil
		return
	}

	getObjConf := func() (icf conf.IConf, err error) {
		icf, err = getLoggerConf()
		if err != nil {
			err = fmt.Errorf("get log conf: %w", err)
			return
		}
		return
	}

	ict, err := conf.NewContainer(getObjConf, compareLoggerConf, newLoggerObj, resetLoggerObj)
	if err != nil {
		err = fmt.Errorf("new conf container: %w", err)
		return
	}

	ct = &LoggerContainer{
		Container: ict,
	}

	return
}

func newLoggerObj(icf conf.IConf) (iobj conf.IObject, err error) {
	cf, ok := icf.(*LoggerConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	iobj, err = NewLogger(cf)
	if err != nil {
		err = fmt.Errorf("new log: %w", err)
		return
	}

	return
}

func compareLoggerConf(iocf, incf conf.IConf) (rst conf.CompareObjConfRst, err error) {
	ocf, ok := iocf.(*LoggerConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	ncf, ok := incf.(*LoggerConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	if *ocf != *ncf {
		rst = conf.CompareObjConfRstNeedReset
		return
	}

	rst = conf.CompareObjConfRstNoNeed

	return
}

func resetLoggerObj(iobj conf.IObject, iocf, incf conf.IConf) (err error) {
	ocf, ok := iocf.(*LoggerConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	ncf, ok := incf.(*LoggerConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	logger, ok := iobj.(*Logger)
	if !ok {
		err = conf.ErrInvalidObjectType
		return
	}

	switch {
	case ncf.Level != ocf.Level:
		if err = logger.SetLevel(ncf); err != nil {
			err = fmt.Errorf("set level: %w", err)
			return
		}
		ocf.Level = ncf.Level
		fallthrough

	case ncf.ReportCaller != ocf.ReportCaller:
		logger.SetReportCaller(ncf)
		ocf.ReportCaller = ncf.ReportCaller
		fallthrough

	case
		ncf.InfLinkName != ocf.InfLinkName,
		ncf.InfRotateHours != ocf.InfRotateHours,
		ncf.InfMaxAgeHours != ocf.InfMaxAgeHours,
		ncf.ErrLinkName != ocf.ErrLinkName,
		ncf.ErrRotateHours != ocf.ErrRotateHours,
		ncf.ErrMaxAgeHours != ocf.ErrMaxAgeHours:

		if err = logger.SetHook(ncf); err != nil {
			err = fmt.Errorf("set hook: %w", err)
			return
		}

		*ocf = *ncf
	default:

	}

	return
}

func (ct *LoggerContainer) MustGetLogger() (logger *Logger) {
	obj := ct.MustGetObj()

	logger, ok := obj.(*Logger)
	if !ok {
		panic(conf.ErrInvalidObjectType)
	}

	return
}

func (ct *LoggerContainer) PutLogger(logger *Logger) {
	ct.PutObj(logger)
}
