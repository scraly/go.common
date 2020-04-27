package server

import (
	"context"
	"net/http"
	"time"

	"github.com/scraly/go.common/pkg/log"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

const (
	prometheusMetricsServerAddress         = "0.0.0.0:9402"
	prometheusMetricsServerShutdownTimeout = 10 * time.Second
	prometheusMetricsServerReadTimeout     = 8 * time.Second
	prometheusMetricsServerWriteTimeout    = 8 * time.Second
	prometheusMetricsServerMaxHeaderBytes  = 1 << 20
)

type prometheusMetricsServer struct {
	http.Server
}

func prometheusHandler() http.Handler {
	return prometheus.Handler()
}

func newPrometheusMetricsServer() *prometheusMetricsServer {

	r := mux.NewRouter()
	r.Handle("/metrics", prometheusHandler())

	// Create server and register prometheus metrics handler
	s := &prometheusMetricsServer{
		Server: http.Server{
			Addr:           prometheusMetricsServerAddress,
			ReadTimeout:    prometheusMetricsServerReadTimeout,
			WriteTimeout:   prometheusMetricsServerWriteTimeout,
			MaxHeaderBytes: prometheusMetricsServerMaxHeaderBytes,
			Handler:        r,
		},
	}

	return s
}

func (s *prometheusMetricsServer) WaitShutdown(ctx context.Context) {
	<-ctx.Done()
	log.For(ctx).Info("Stopping Prometheus metrics server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, prometheusMetricsServerShutdownTimeout)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.For(shutdownCtx).Error("Prometheus metrics server shutdown error", zap.Error(err))
		return
	}

	log.For(ctx).Info("Prometheus metrics server gracefully stopped")
}

// StartPrometheusMetricsServer is used to start the prometheus server
func StartPrometheusMetricsServer(ctx context.Context) {
	s := newPrometheusMetricsServer()

	go func() {

		log.For(ctx).Info("Metric server listening ...", zap.String("addr", s.Addr))
		if err := s.ListenAndServe(); err != nil {
			log.For(ctx).Error("Error running prometheus metrics server", zap.Error(err))
			return
		}

		glog.Infof("Prometheus metrics server exited")

	}()

	s.WaitShutdown(ctx)
}
