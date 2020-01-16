package wrapper

import (
	"context"
	"fmt"
	"github.com/Gitforxuyang/microBase/trace"
	"github.com/Gitforxuyang/microBase/util"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
)

//做请求时的链路日志
func NewCallTraceWrapper(ot opentracing.Tracer) client.CallWrapper {
	return func(callFunc client.CallFunc) client.CallFunc {
		return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
			ctx, span, err := trace.StartSpanFromContext(ctx, ot, fmt.Sprintf("%s", req.Method()))
			span.SetTag("span.kind", "client")
			if err != nil {
				util.Error(ctx, err.Error())
				return err
			}
			defer span.Finish()
			err = callFunc(ctx, node, req, rsp, opts)
			if err != nil {
				util.ErrorKV(ctx, logrus.Fields{"req": req.Body(), "rsp": rsp}, err.Error())
				ext.Error.Set(span, true)
				span.LogKV("error.kind", err.Error(), "message", err.Error())
			}
			return err
		}
	}
}

func NewClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		//设置重试次数=0
		if err := c.Init(client.Retries(0)); err != nil {
			util.MicroLog.Error(err)
		}
		return c
	}
}
