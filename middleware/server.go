package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/util/log"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"gitlab.neoclub.cn/cms/go/microbase"
	"gitlab.neoclub.cn/cms/go/microbase/trace"
)

func NewLogWrapper() server.HandlerWrapper {
	return func(handlerFunc server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) (err error) {
			defer func() {
				r := recover()
				//TODO:sentry异常捕获
				if r != nil {
					microBase.InfoKV(ctx,
						logrus.Fields{"service": req.Service(),
							"endpoint": req.Endpoint(), "method": req.Method(), "body": req.Body()}, fmt.Sprintf("panic: %s", r))
					err = errors.New(fmt.Sprintf("panic: %s", r))
				}
			}()
			//进入时打印日志
			microBase.InfoKV(ctx,
				logrus.Fields{"service": req.Service(),
					"endpoint": req.Endpoint(), "method": req.Method(), "body": req.Body()}, "")
			err = handlerFunc(ctx, req, rsp)
			//退出时打印日志
			microBase.Info(ctx, rsp)
			if err != nil {
				log.Log(err.Error())
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
					log.Log(fmt.Sprintf("panic: %s", r))
				}
			}()
			name := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
			ctx, span, err := trace.StartSpanFromContext(ctx, ot, name)
			s, ok := span.Context().(jaeger.SpanContext)
			if !ok {
				microBase.Info(ctx, "spanContext转化失败")
			} else {
				//如果转化正常，将traceId携带到ctx上
				traceId := s.TraceID().String()
				ctx = context.WithValue(ctx, traceId, traceId)
			}
			if err != nil {
				return err
			}
			defer span.Finish()
			err = h(ctx, req, rsp)
			//当处理函数返回错误时，记录进链路
			if err != nil {
				ext.Error.Set(span, true)
				span.LogKV(err)
			}
			return err
		}
	}
}
