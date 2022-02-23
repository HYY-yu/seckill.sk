// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package api

import (
	"github.com/HYY-yu/seckill.sk/internal/pkg/cache"
	"github.com/HYY-yu/seckill.sk/internal/pkg/db"
	"github.com/HYY-yu/seckill.sk/internal/service/sk/api/handler"
	"github.com/HYY-yu/seckill.sk/internal/service/sk/api/repo"
	"github.com/HYY-yu/seckill.sk/internal/service/sk/api/svc"
)

// Injectors from wire.go:

// initHandlers init Handlers.
func initHandlers(d db.Repo, c cache.Repo) (*Handlers, error) {
	skRepo := repo.NewSKRepo()
	skSvc := svc.NewSKSvc(d, c, skRepo)
	skHandler := handler.NewSKHandler(skSvc)
	handlers := NewHandlers(skHandler)
	return handlers, nil
}
