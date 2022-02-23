package handler

import (
	"github.com/HYY-yu/seckill.sk/internal/service/sk/api/svc"
)

type SKHandler struct {
	skSvc *svc.SKSvc
}

func NewSKHandler(skSvc *svc.SKSvc) *SKHandler {
	return &SKHandler{
		skSvc: skSvc,
	}
}
