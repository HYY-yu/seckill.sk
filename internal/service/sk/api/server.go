package api

import (
	"errors"

	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"

	"github.com/HYY-yu/seckill.pkg/pkg/metrics"

	"github.com/HYY-yu/seckill.sk/internal/pkg/middleware"
	"github.com/HYY-yu/seckill.sk/internal/service/sk/api/handler"
	"github.com/HYY-yu/seckill.sk/internal/service/sk/config"

	"github.com/HYY-yu/seckill.pkg/pkg/jaeger"

	"github.com/HYY-yu/seckill.sk/internal/pkg/cache"
	"github.com/HYY-yu/seckill.sk/internal/pkg/core"
	"github.com/HYY-yu/seckill.sk/internal/pkg/db"
)

type Handlers struct {
	skHandler *handler.SKHandler
}

func NewHandlers(
	skHandler *handler.SKHandler,
) *Handlers {
	return &Handlers{
		skHandler: skHandler,
	}
}

type Server struct {
	Logger  *zap.Logger
	Engine  core.Engine
	DB      db.Repo
	Cache   cache.Repo
	Trace   *trace.TracerProvider
	Middles middleware.Middleware
}

func NewApiServer(logger *zap.Logger) (*Server, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}
	s := &Server{}
	s.Logger = logger

	dbRepo, err := db.New()
	if err != nil {
		logger.Fatal("new db err", zap.Error(err))
	}
	s.DB = dbRepo

	cacheRepo, err := cache.New()
	if err != nil {
		logger.Fatal("new cache err", zap.Error(err))
	}
	s.Cache = cacheRepo

	// Jaeger
	tp, err := jaeger.InitJaeger(config.Get().Server.ServerName, config.Get().Jaeger.UdpEndpoint)
	if err != nil {
		logger.Error("jaeger error", zap.Error(err))
	}
	s.Trace = tp

	// Metrics
	metrics.InitMetrics(config.Get().Server.ServerName, "api")

	opts := make([]core.Option, 0)
	opts = append(opts, core.WithEnableCors())
	opts = append(opts, core.WithRecordMetrics(metrics.RecordMetrics))
	if !config.Get().Server.Pprof {
		opts = append(opts, core.WithDisablePProf())
	}

	engine, err := core.New(logger, opts...)
	if err != nil {
		panic(err)
	}
	s.Engine = engine

	s.Middles = middleware.New(logger)

	// Init Repo Svc Handler
	c, err := initHandlers(s.DB, s.Cache)
	if err != nil {
		panic(err)
	}

	s.Route(c)
	return s, nil
}