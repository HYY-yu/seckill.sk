package svc

import (
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

type SKSvc struct {
	DB    db.Repo
	Cache cache_v2.Repo

	SKRepo     repo.SKRepo
	ShopClient proto.ShopClient
}

func NewSKSvc(
	db db.Repo,
	ca cache_v2.Repo,
	goodsRepo repo.SKRepo,
	shopClient proto.ShopClient,
) *SKSvc {
	return &SKSvc{
		DB:         db,
		Cache:      ca,
		SKRepo:     goodsRepo,
		ShopClient: shopClient,
	}
}

func (s *SKSvc) List(sctx core.SvcContext, pr *page.PageRequest) (*page.Page, error) {
	ctx := sctx.Context()
	mgr := s.SKRepo.Mgr(ctx, s.DB.GetDb(ctx))

	limit, offset := pr.GetLimitAndOffset()
	pr.AddAllowSortField(model.SecKillColumns.CreateTime)
	pr.AddAllowSortField(model.SecKillColumns.StartTime)
	sort, _ := pr.Sort()

	data, err := mgr.ListSK(limit, offset, pr.Filter, sort)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	count, err := mgr.CountSK(pr.Filter)
	if err != nil {
		return nil, response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	var result = make([]model.SKListResp, len(data))
	// 查询shop
	shopIds := make([]int64, 0, len(data))
	for _, e := range data {
		if e.ShopID == 0 {
			continue
		}
		shopIds = append(shopIds, int64(e.ShopID))
	}
	shopResp, err := s.ShopClient.List(ctx, &proto.ListReq{
		PageNo:    1,
		PageSize:  int32(limit),
		SortBy:    "",
		FieldList: []string{"name", "desc", "count", "create_time"},
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
		// 查询
		sdata := new(proto.ShopData)
		for _, d := range shopDatas {
			if d.Id == int64(e.ShopID) {
				sdata = d
			}
		}

		r := model.SKListResp{
			ID:     e.ID,
			ShopId: e.ShopID,
			ShopInfo: struct {
				ShopName  string `json:"shop_name"`
				ShopDesc  string `json:"shop_desc"`
				ShopCount int    `json:"shop_count"`
			}{
				ShopName:  sdata.GetName(),
				ShopDesc:  sdata.GetDesc(),
				ShopCount: int(sdata.GetCount()),
			},
			StartTime:  e.StartTime,
			EndTime:    e.EndTime,
			Status:     int(e.Status),
			CreateTime: e.CreateTime,
		}
		result[i] = r
	}
	return page.NewPage(
		count,
		result,
	), nil
}

func (s *SKSvc) Add(sctx core.SvcContext, param *model.SKAdd) error {
	// 确保 ShopId 存在
	ctx := sctx.Context()
	mgr := s.SKRepo.Mgr(ctx, s.DB.GetDb(ctx))

	ok, err := mgr.WithOptions(mgr.WithShopID(param.ShopId)).HasRecord()
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	if !ok {
		return response.NewError(
			http.StatusBadRequest,
			response.ParamBindError,
			"请输入正确的ShopId",
		)
	}

	startTime := time.Unix(int64(param.StartTime), 0)
	endTime := time.Unix(int64(param.EndTime), 0)

	if time.Until(startTime) < time.Second*30 {
		return response.NewError(
			http.StatusBadRequest,
			response.ParamBindError,
			"start_time 过近",
		)
	}

	// 存数据库（事务）
	tx := s.DB.GetDb(ctx)

}
