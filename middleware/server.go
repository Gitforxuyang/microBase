package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gitforxuyang/microBase/trace"
	"github.com/Gitforxuyang/microBase/util"
	"github.com/micro/go-micro/server"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
)

func NewLogWrapper() server.HandlerWrapper {
	return func(handlerFunc server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) (err error) {
			defer func() {
				r := recover()
				//TODO:sentry异常捕获
				if r != nil {
					util.Error(ctx, fmt.Sprintf("panic: %s", r))
					err = errors.New(fmt.Sprintf("panic: %s", r))
				}
			}()
			//进入时打印日志
			util.InfoKV(ctx,
				logrus.Fields{"method": req.Method(), "body": req.Body()}, "")
			err = handlerFunc(ctx, req, rsp)
			//退出时打印日志
			util.Info(ctx, rsp)
			if err != nil {
				util.Error(ctx, err.Error())
			}
			return err
		}
	}
}

func NewTraceWrapper(ot opentracing.Tracer) server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			defer func() {
				r := recover()
				if r != nil {
					util.Error(ctx, fmt.Sprintf("panic: %s", r))
				}
			}()
			name := req.Method()
			ctx, span, err := trace.StartSpanFromContext(ctx, ot, name)
			span.SetTag("span.kind", "server")
			s, ok := span.Context().(jaeger.SpanContext)
			if !ok {
				util.Info(ctx, "spanContext转化失败")
			} else {
				//如果转化正常，将traceId携带到ctx上
				traceId := s.TraceID().String()
				ctx = context.WithValue(ctx, "traceId", traceId)
			}
			if err != nil {
				return err
			}
			defer span.Finish()
			err = h(ctx, req, rsp)
			//当处理函数返回错误时，记录进链路
			if err != nil {
				ext.Error.Set(span, true)
				span.LogKV("error.kind", err.Error(), "message", err.Error())
			}
			return err
		}
	}
}
