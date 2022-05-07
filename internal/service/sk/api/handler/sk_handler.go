package handler

import (
	"context"
	"net/http"

	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/pkg/page"
	"github.com/HYY-yu/seckill.pkg/pkg/response"
	"github.com/gogf/gf/v2/util/gvalid"

	"github.com/HYY-yu/seckill.sk/internal/service/sk/api/svc"
	"github.com/HYY-yu/seckill.sk/internal/service/sk/model"
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

func (s *SKHandler) Add(c core.Context) {
	params := &model.SKAdd{}
	err := c.ShouldBindForm(params)
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
		return
	}

	validErr := gvalid.CheckStruct(context.Background(), params, nil)
	if validErr != nil {
		c.AbortWithError(response.NewError(
			http.StatusBadRequest,
			response.ParamBindError,
			validErr.Error(),
		))
		return
	}

	err = s.skSvc.Add(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}

func (s *SKHandler) Delete(c core.Context) {
	type DeleteParam struct {
		Id int `form:"id" v:"required"`
	}
	param := &DeleteParam{}
	err := c.ShouldBindForm(param)
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
		return
	}

	validErr := gvalid.CheckStruct(context.Background(), param, nil)
	if validErr != nil {
		c.AbortWithError(response.NewError(
			http.StatusBadRequest,
			response.ParamBindError,
			validErr.Error(),
		))
		return
	}

	err = s.skSvc.Delete(c.SvcContext(), param.Id)
	c.AbortWithError(err)
	c.Payload(nil)
}
