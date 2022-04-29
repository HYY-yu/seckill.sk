//go:build wireinject
// +build wireinject

//go:generate wire gen .
package api

import (
	"github.com/HYY-yu/seckill.shop/proto"
	"github.com/google/wire"

	"github.com/HYY-yu/seckill.pkg/cache_v2"
	"github.com/HYY-yu/seckill.pkg/db"

	"github.com/HYY-yu/seckill.sk/internal/service/sk/api/handler"
	"github.com/HYY-yu/seckill.sk/internal/service/sk/api/repo"
	"github.com/HYY-yu/seckill.sk/internal/service/sk/api/svc"
)

// initHandlers init Handlers.
func initHandlers(d db.Repo, c cache_v2.Repo, client proto.ShopClient) (*Handlers, error) {
	panic(wire.Build(repo.NewSKRepo, svc.NewSKSvc, handler.NewSKHandler, NewHandlers))
}
