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

type OrderHandler struct {
	orderSvc *svc.OrderSvc
}

func NewOrderHandler(orderSvc *svc.OrderSvc) *OrderHandler {
	return &OrderHandler{
		orderSvc: orderSvc,
	}
}

func (o *OrderHandler) List(c core.Context) {
	err := c.RequestContext().Request.ParseForm()
	if err != nil {
		c.AbortWithError(response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		).WithErr(err))
		return
	}
	pageRequest := page.NewPageFromRequest(c.RequestContext().Request.Form)

	data, err := o.orderSvc.List(c.SvcContext(), pageRequest)
	c.AbortWithError(err)
	c.Payload(data)
}

func (o *OrderHandler) Join(c core.Context) {
	params := &model.OrderJoin{}

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

	err = o.orderSvc.Join(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(nil)
}
