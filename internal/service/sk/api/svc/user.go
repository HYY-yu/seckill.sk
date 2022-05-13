package svc

import (
	"net/http"

	"github.com/HYY-yu/seckill.pkg/cache_v2"
	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/HYY-yu/seckill.pkg/pkg/login"
	"github.com/HYY-yu/seckill.pkg/pkg/login/model"
	"github.com/HYY-yu/seckill.pkg/pkg/response"

	"github.com/HYY-yu/seckill.sk/internal/service/sk/api/repo"
	"github.com/HYY-yu/seckill.sk/internal/service/sk/config"
)

type UserSvc struct {
	DB    db.Repo
	Cache cache_v2.Repo

	UserRepo repo.UserRepo

	system login.LoginTokenSystem
}

func NewUserSvc(
	db db.Repo,
	ca cache_v2.Repo,
	userRepo repo.UserRepo,
) *UserSvc {
	svc := &UserSvc{
		DB:       db,
		Cache:    ca,
		UserRepo: userRepo,
	}
	jwtCfg := config.Get().JWT

	switch jwtCfg.Type {
	case "refresh_token":
		cfg := &login.RefreshTokenConfig{
			Secret:          jwtCfg.Secret,
			ExpireDuration:  jwtCfg.ExpireDuration,
			RefreshDuration: jwtCfg.RefreshDuration,
		}

		svc.system = login.NewByRefreshToken(cfg, ca)
	case "black_list":
		cfg := &login.BlackListConfig{
			Secret:         jwtCfg.Secret,
			ExpireDuration: jwtCfg.ExpireDuration,
		}
		svc.system = login.NewByBlackList(cfg, ca)
	}
	return svc
}

func (u *UserSvc) Login(sctx core.SvcContext, param *model.LoginParams) (*model.LoginResponse, error) {
	ctx := sctx.Context()
	mgr := u.UserRepo.Mgr(ctx, u.DB.GetDb(ctx))

	// 查询此用户
	user, err := mgr.WithOptions(mgr.WithUserName(param.UserName)).Get()
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	if user.ID == 0 {
		// 注册
		user.UserName = param.UserName
		err = mgr.CreateUsers(&user)
		if err != nil {
			return nil, response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
	}

	// 派发Token
	resp, err := u.system.GenerateToken(ctx, user.ID, user.UserName)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	return resp, err
}

func (u *UserSvc) RefreshToken(sctx core.SvcContext, oldToken string) (*model.LoginResponse, error) {
	ctx := sctx.Context()

	resp, err := u.system.RefreshToken(ctx, oldToken)
	if err != nil {
		return nil, response.NewError(
			http.StatusInternalServerError,
			response.ServerError,
			"请检查RefreshToken是否正确",
		).WithErr(err)
	}

	return resp, err
}

func (u *UserSvc) Unlogin(sctx core.SvcContext) error {
	ctx := sctx.Context()

	err := u.system.TokenCancelById(ctx, int(sctx.UserId()), sctx.UserName())
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}

func (u *UserSvc) CheckBlackList(sctx core.SvcContext, accessToken string) (bool, error) {
	ctx := sctx.Context()

	blackListSystem, ok := u.system.(*login.BlackListSystem)
	if !ok {
		return false, nil
	}

	return blackListSystem.CheckBlackList(ctx, accessToken)
}
