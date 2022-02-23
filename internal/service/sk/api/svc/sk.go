package svc

import (
	"github.com/HYY-yu/seckill.sk/internal/pkg/cache"
	"github.com/HYY-yu/seckill.sk/internal/pkg/db"
	"github.com/HYY-yu/seckill.sk/internal/service/sk/api/repo"
)

type SKSvc struct {
	DB    db.Repo
	Cache cache.Repo

	SKRepo repo.SKRepo
}

func NewSKSvc(db db.Repo, ca cache.Repo, goodsRepo repo.SKRepo) *SKSvc {
	return &SKSvc{
		DB:     db,
		Cache:  ca,
		SKRepo: goodsRepo,
	}
}
