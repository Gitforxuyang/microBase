package trace

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport"
	"io"
	"time"
)

func NewTracer(serviceName string, addr string) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeProbabilistic,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	sender := transport.NewHTTPTransport(addr)
	reporter := jaeger.NewRemoteReporter(sender)
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		config.Reporter(reporter),
	)
	return tracer, closer, err
}
