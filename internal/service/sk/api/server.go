package api

import (
	"errors"

	"github.com/HYY-yu/seckill.pkg/cache_v2"
	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/core/middleware"
	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/HYY-yu/seckill.shop/proto/client"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"

	"github.com/HYY-yu/seckill.pkg/pkg/metrics"

	"github.com/HYY-yu/seckill.sk/internal/service/sk/api/handler"
	"github.com/HYY-yu/seckill.sk/internal/service/sk/config"

	"github.com/HYY-yu/seckill.pkg/pkg/jaeger"
)

type Handlers struct {
	skHandler    *handler.SKHandler
	loginHandler *handler.LoginHandler
	orderHandler *handler.OrderHandler
}

func NewHandlers(
	skHandler *handler.SKHandler,
	loginHandler *handler.LoginHandler,
	orderHandler *handler.OrderHandler,
) *Handlers {
	return &Handlers{
		skHandler:    skHandler,
		loginHandler: loginHandler,
		orderHandler: orderHandler,
	}
}

type Server struct {
	Logger  *zap.Logger
	Engine  core.Engine
	DB      db.Repo
	Cache   cache_v2.Repo
	Trace   *trace.TracerProvider
	Middles middleware.Middleware
}

func NewApiServer(logger *zap.Logger) (*Server, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}
	s := &Server{}
	s.Logger = logger
	cfg := config.Get()

	dbRepo, err := db.New(&db.DBConfig{
		User:            cfg.MySQL.Base.User,
		Pass:            cfg.MySQL.Base.Pass,
		Addr:            cfg.MySQL.Base.Addr,
		Name:            cfg.MySQL.Base.Name,
		MaxOpenConn:     cfg.MySQL.Base.MaxOpenConn,
		MaxIdleConn:     cfg.MySQL.Base.MaxIdleConn,
		ConnMaxLifeTime: cfg.MySQL.Base.ConnMaxLifeTime,
		ServerName:      cfg.Server.ServerName,
	})
	if err != nil {
		logger.Fatal("new db err", zap.Error(err))
	}
	s.DB = dbRepo

	cacheRepo, err := cache_v2.New(cfg.Server.ServerName, &cache_v2.RedisConf{
		Addr:         cfg.Redis.Addr,
		Pass:         cfg.Redis.Pass,
		Db:           cfg.Redis.Db,
		MaxRetries:   cfg.Redis.MaxRetries,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})
	if err != nil {
		logger.Fatal("new cache err", zap.Error(err))
	}
	s.Cache = cacheRepo

	// Jaeger
	var tp *trace.TracerProvider
	if cfg.Jaeger.StdOut {
		tp, err = jaeger.InitStdOutForDevelopment(cfg.Server.ServerName, cfg.Jaeger.UdpEndpoint)
	} else {
		tp, err = jaeger.InitJaeger(cfg.Server.ServerName, cfg.Jaeger.UdpEndpoint)
	}
	if err != nil {
		logger.Error("jaeger error", zap.Error(err))
	}
	s.Trace = tp

	// Metrics
	metrics.InitMetrics(cfg.Server.ServerName, "api")

	opts := make([]core.Option, 0)
	opts = append(opts, core.WithEnableCors())
	opts = append(opts, core.WithRecordMetrics(metrics.RecordMetrics))
	if !config.Get().Server.Pprof {
		opts = append(opts, core.WithDisablePProf())
	}

	engine, err := core.New(cfg.Server.ServerName, logger, opts...)
	if err != nil {
		panic(err)
	}
	s.Engine = engine

	s.Middles = middleware.New(logger, cfg.JWT.Secret)

	// GRPC Client
	shopClient, err := client.Connect(cfg.Server.Grpc.ShopHost)
	if err != nil {
		panic(err)
	}

	// Init Repo Svc Handler
	c, err := initHandlers(s.DB, s.Cache, shopClient)
	if err != nil {
		panic(err)
	}

	s.Route(c)
	return s, nil
}
