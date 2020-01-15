package util

import (
	"context"
	"github.com/micro/go-micro/util/log"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	microLog *logrus.Logger
)

func InitLog() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	logger := logrus.New()
	baseLog := &baseLog{logger: logger}
	log.SetLogger(baseLog)
	microLog = logger
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
