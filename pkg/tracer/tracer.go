/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package tracer

import (
	"fmt"
	"time"

	"github.com/scraly/go.common/pkg/log"
	"go.uber.org/zap"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
)

// Tracer instance
var tracer opentracing.Tracer

// SetTracer can be used by unit tests to provide a NoopTracer instance. Real users should always
// use the InitTracing func.
func SetTracer(initializedTracer opentracing.Tracer) {
	tracer = initializedTracer
}

// Sampling Server URL
var samplingServerURL string

// SetSamplingServerURL defines the global server url
func SetSamplingServerURL(url string) {
	samplingServerURL = url
}

// InitTracing connects the calling service to Zipkin and initializes the tracer.
func InitTracing(serviceName string, logger log.LoggerFactory) opentracing.Tracer {
	logger.Bg().Debug("Initializing tracer", zap.String("svc", serviceName))

	// cfg := config.Configuration{
	// 	Reporter: &config.ReporterConfig{
	// 		LocalAgentHostPort: "127.0.0.1:6831",
	// 	},
	// }

	// tracer, _, err := cfg.New(
	// 	serviceName,
	// 	config.Logger(jaegerLoggerAdapter{logger.Bg()}),
	// 	config.ZipkinSharedRPCSpan(true),
	// 	config.Gen128Bit(true),
	// )

	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	injector := jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator)
	extractor := jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator)

	// Zipkin shares span ID between client and server spans; it must be enabled via the following option.
	zipkinSharedRPCSpan := jaeger.TracerOptions.ZipkinSharedRPCSpan(true)

	// sender, err := jaeger.NewUDPTransport("jaeger-agent.istio-system:5775", 0)
	sender, err := jaeger.NewUDPTransport(samplingServerURL, 0)
	if err != nil {
		logger.Bg().Fatal("cannot initialize Jaeger Tracer", zap.Error(err))
	}

	tracer, _ := jaeger.NewTracer(
		serviceName,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(
			sender,
			jaeger.ReporterOptions.BufferFlushInterval(1*time.Second)),
		injector,
		extractor,
		zipkinSharedRPCSpan,
		jaeger.TracerOptions.Logger(jaegerLoggerAdapter{logger.Bg()}),
	)

	// if err != nil {
	// 	logger.Bg().Fatal("cannot initialize Jaeger Tracer", zap.Error(err))
	// }

	return tracer
}

type jaegerLoggerAdapter struct {
	logger log.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}
