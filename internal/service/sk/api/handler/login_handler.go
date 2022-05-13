package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/pkg/login/model"
	"github.com/HYY-yu/seckill.pkg/pkg/response"
	"github.com/gogf/gf/v2/util/gvalid"

	"github.com/HYY-yu/seckill.sk/internal/service/sk/api/svc"
	"github.com/HYY-yu/seckill.sk/internal/service/sk/config"
)

type LoginHandler struct {
	userSvc *svc.UserSvc
}

func NewLoginHandler(userSvc *svc.UserSvc) *LoginHandler {
	return &LoginHandler{
		userSvc: userSvc,
	}
}

func (s *LoginHandler) Login(c core.Context) {
	params := &model.LoginParams{}
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

	data, err := s.userSvc.Login(c.SvcContext(), params)
	c.AbortWithError(err)
	c.Payload(data)
}

func (s *LoginHandler) Unlogin(c core.Context) {
	err := s.userSvc.Unlogin(c.SvcContext())
	c.AbortWithError(err)
}

func (s *LoginHandler) RefreshToken(c core.Context) {
	type Param struct {
		RefreshToken string `form:"refresh_token" v:"required"`
	}
	params := &Param{}
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

	data, err := s.userSvc.RefreshToken(c.SvcContext(), params.RefreshToken)
	c.AbortWithError(err)
	c.Payload(data)
}

func (s *LoginHandler) getAccessToken(c core.Context) string {
	auth := c.GetHeader("Authorization")
	stripBearer := func(tok string) string {
		// Should be a bearer token
		if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
			return tok[7:]
		}
		return tok
	}
	auth = stripBearer(auth)
	return auth
}

func (s *LoginHandler) CheckBlackList(c core.Context) {
	cfg := config.Get().JWT

	if cfg.Type != "black_list" {
		c.RequestContext().Next()
		return
	}

	// 读取AccessToken
	auth := s.getAccessToken(c)
	has, err := s.userSvc.CheckBlackList(c.SvcContext(), auth)
	if err != nil {
		// 读不到，报错，这里不管
		c.Logger().Error("can't check black list, skip . ")
		c.RequestContext().Next()
		return
	}

	if has {
		err = response.NewErrorAutoMsg(
			http.StatusUnauthorized,
			response.AuthorizationError,
		).WithErr(errors.New("Header 中缺少 Authorization 参数 "))
		c.AbortWithError(err)
		return
	}
}
