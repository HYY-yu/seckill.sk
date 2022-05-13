package repo

import (
	"context"

	"gorm.io/gorm"
)

type UserRepo interface {
	Mgr(ctx context.Context, db *gorm.DB) *_UsersMgr
}

type userRepo struct {
}

func NewUserRepo() UserRepo {
	return &userRepo{}
}

func (*userRepo) Mgr(ctx context.Context, db *gorm.DB) *_UsersMgr {
	mgr := UsersMgr(ctx, db)
	return mgr
}

// ------- 自定义方法 -------
