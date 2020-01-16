package util

import (
	"context"
	"github.com/micro/go-micro/util/log"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	microLog *logrus.Logger
	MicroLog *logrus.Logger
)

func InitLog() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	baseLog := &baseLog{logger: logger}
	log.SetLogger(baseLog)
	microLog = logger
	MicroLog = logger
}

type baseLog struct {
	logger *logrus.Logger
}

func (b *baseLog) Log(v ...interface{}) {
	b.logger.Info(v)
}
func (b *baseLog) Logf(format string, v ...interface{}) {
	b.logger.Infof(format, v)
}

//根据ctx获取traceId 让全局的打印日志都含有traceId
func Info(ctx context.Context, msg interface{}) {
	microLog.WithFields(logrus.Fields{"traceId": ctx.Value("traceId")}).Info(msg)
}

//打印日志时携带对象
func InfoKV(ctx context.Context, fields logrus.Fields, msg interface{}) {
	fields["traceId"] = ctx.Value("traceId")
	microLog.WithFields(fields).Info(msg)
}

//根据ctx获取traceId 让全局的打印日志都含有traceId
func Error(ctx context.Context, msg interface{}) {
	microLog.WithFields(logrus.Fields{"traceId": ctx.Value("traceId")}).Error(msg)
}

//打印日志时携带对象
func ErrorKV(ctx context.Context, fields logrus.Fields, msg interface{}) {
	fields["traceId"] = ctx.Value("traceId")
	microLog.WithFields(fields).Error(msg)
}

//根据ctx获取traceId 让全局的打印日志都含有traceId
func Warn(ctx context.Context, msg interface{}) {
	microLog.WithFields(logrus.Fields{"traceId": ctx.Value("traceId")}).Warn(msg)
}

//打印日志时携带对象
func WarnKV(ctx context.Context, fields logrus.Fields, msg interface{}) {
	fields["traceId"] = ctx.Value("traceId")
	microLog.WithFields(fields).Warn(msg)
}
