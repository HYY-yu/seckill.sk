package svc

import (
	"errors"
	"net/http"
	"time"

	"github.com/HYY-yu/seckill.pkg/cache_v2"
	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/HYY-yu/seckill.pkg/pkg/page"
	"github.com/HYY-yu/seckill.pkg/pkg/response"
	"github.com/HYY-yu/seckill.shop/proto"

	"github.com/HYY-yu/seckill.sk/internal/service/sk/api/repo"
	"github.com/HYY-yu/seckill.sk/internal/service/sk/model"
)

type OrderSvc struct {
	DB    db.Repo
	Cache cache_v2.Repo

	SKRepo     repo.SKRepo
	OrderRepo  repo.OrderRepo
	UserRepo   repo.UserRepo
	ShopClient proto.ShopClient
}

func NewOrderSvc(
	db db.Repo,
	ca cache_v2.Repo,
	skRepo repo.SKRepo,
	orderRepo repo.OrderRepo,
	userRepo repo.UserRepo,
	shopClient proto.ShopClient,
) *OrderSvc {
	svc := &OrderSvc{
		DB:         db,
		Cache:      ca,
		SKRepo:     skRepo,
		OrderRepo:  orderRepo,
		UserRepo:   userRepo,
		ShopClient: shopClient,
	}

	return svc
}

func (o *OrderSvc) List(sctx core.SvcContext, pr *page.PageRequest) (*page.Page, error) {
	ctx := sctx.Context()
	mgr := o.OrderRepo.Mgr(ctx, o.DB.GetDb(ctx))
	userMgr := o.UserRepo.Mgr(ctx, o.DB.GetDb(ctx))

	limit, offset := pr.GetLimitAndOffset()
	pr.AddAllowSortField(model.OrderColumns.CreateTime)
	sort, _ := pr.Sort()

	data, err := mgr.ListOrder(limit, offset, pr.Filter, sort)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	count, err := mgr.CountOrder(pr.Filter)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	var result = make([]model.OrderListResp, len(data))
	// 查询shop
	shopIds := make([]int64, 0, len(data))
	for _, e := range data {
		if e.ShopID == 0 {
			continue
		}
		shopIds = append(shopIds, int64(e.ShopID))
	}
	shopResp, err := o.ShopClient.List(ctx, &proto.ListReq{
		PageNo:    1,
		PageSize:  int32(limit),
		SortBy:    "",
		FieldList: []string{"name"},
		ShopIds:   shopIds,
	})
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	shopDatas := shopResp.GetData()

	for i, e := range data {
		sdata := new(proto.ShopData)
		for _, d := range shopDatas {
			if d.Id == int64(e.ShopID) {
				sdata = d
			}
		}

		user, err := userMgr.WithOptions(userMgr.WithID(e.UserID)).Get()
		if err != nil {
			return nil, response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}

		r := model.OrderListResp{
			ID:         e.ID,
			ShopId:     e.ShopID,
			ShopName:   sdata.GetName(),
			UserId:     e.UserID,
			UserName:   user.UserName,
			SkID:       e.SecID,
			CreateTime: e.CreateTime,
		}

		result[i] = r
	}
	return page.NewPage(
		count,
		result,
	), nil
}

// Join
// 参与秒杀
// 0. 校验秒杀活动、校验商品数量
// 1. 校验用户是否购买过此商品
// 2. 加入订单表
// 3. 扣减商品库存
func (o *OrderSvc) Join(sctx core.SvcContext, param *model.OrderJoin) error {
	ctx := sctx.Context()
	userId := sctx.UserId()

	mgr := o.OrderRepo.Mgr(ctx, o.DB.GetDb(ctx))
	skMgr := o.SKRepo.Mgr(ctx, o.DB.GetDb(ctx))

	sk, err := skMgr.WithOptions(skMgr.WithID(param.SKId)).
		WithSelects(
			model.SecKillColumns.ID,
			model.SecKillColumns.Status,
			model.SecKillColumns.ShopID,
		).
		Get()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	if sk.Status != model.SKShopping {
		// return
		return response.NewErrorWithStatusOk(
			10010,
			"此活动暂不支持参与",
		)
	}
	has, err := mgr.WithOptions(mgr.WithShopID(sk.ShopID), mgr.WithUserID(int(userId))).HasRecord()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	if has {
		// 此用户购买过
		return response.NewErrorWithStatusOk(
			10011,
			"您已参与过此活动",
		)
	}

	// 商品数量
	shopResp, err := o.ShopClient.List(ctx, &proto.ListReq{
		PageNo:    1,
		PageSize:  1,
		SortBy:    "",
		FieldList: []string{"name", "count"},
		ShopIds:   []int64{int64(sk.ShopID)},
	})
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	if len(shopResp.GetData()) <= 0 {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(errors.New("不存在此商品! "))
	}

	shop := shopResp.GetData()[0]
	shopCount := shop.GetCount()

	if shopCount < 0 {
		// 商品售罄
		err = skMgr.UpdateSecKill(&model.SecKill{
			ID:      sk.ID,
			Status:  model.SKFinish,
			EndTime: int(time.Now().Unix()),
		})
		if err != nil {
			return response.NewErrorAutoMsg(
				http.StatusInternalServerError,
				response.ServerError,
			).WithErr(err)
		}
	}

	tx := o.DB.GetDb(ctx).Begin()
	mgr.UpdateDB(tx)

	// 插入订单
	err = mgr.CreateOrder(&model.Order{
		SecID:      sk.ID,
		ShopID:     sk.ShopID,
		UserID:     int(userId),
		CreateTime: int(time.Now().Unix()),
	})
	if err != nil {
		tx.Rollback()
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	_, err = o.ShopClient.Incr(ctx, &proto.IncrReq{
		N:      -1,
		ShopId: int64(sk.ShopID),
	})

	if err != nil {
		tx.Rollback()
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	tx.Commit()
	mgr.Reset()
	return nil
}
