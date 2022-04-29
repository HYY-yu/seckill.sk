package handler

import (
	"net/http"

	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/pkg/page"
	"github.com/HYY-yu/seckill.pkg/pkg/response"

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

func (s *SKHandler) List(c core.Context) {
	err := c.RequestContext().Request.ParseForm()
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
		return
	}
	pageRequest := page.NewPageFromRequest(c.RequestContext().Request.Form)

	data, err := s.skSvc.List(c.SvcContext(), pageRequest)
	c.AbortWithError(err)
	c.Payload(data)
}
