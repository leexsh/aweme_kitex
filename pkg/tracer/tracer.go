package tracer

import (
	"aweme_kitex/pkg/logger"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func InitJaeger(service string) {

	// cfg, _ := jae
	cfg, _ := jaegercfg.FromEnv()
	cfg.ServiceName = service
	tracer, _, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		logger.Error(err)
	}
	opentracing.SetGlobalTracer(tracer)
}
