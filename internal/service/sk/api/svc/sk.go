package svc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/HYY-yu/seckill.pkg/cache"
	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/HYY-yu/seckill.pkg/pkg/elastic_job"
	"github.com/HYY-yu/seckill.pkg/pkg/encrypt"
	"github.com/HYY-yu/seckill.pkg/pkg/page"
	"github.com/HYY-yu/seckill.pkg/pkg/response"
	"github.com/HYY-yu/seckill.shop/proto"
	"github.com/google/uuid"
	"github.com/spf13/cast"

	"github.com/HYY-yu/seckill.sk/internal/service/sk/api/repo"
	"github.com/HYY-yu/seckill.sk/internal/service/sk/model"
)

type SKSvc struct {
	DB    db.Repo
	Cache cache.Repo

	SKRepo     repo.SKRepo
	ShopClient proto.ShopClient
}

func NewSKSvc(
	db db.Repo,
	ca cache.Repo,
	goodsRepo repo.SKRepo,
	shopClient proto.ShopClient,
) *SKSvc {
	svc := &SKSvc{
		DB:         db,
		Cache:      ca,
		SKRepo:     goodsRepo,
		ShopClient: shopClient,
	}
	// 延时任务注册
	elastic_job.Get().RegisterHandler(model.SKDelayAddTag, svc.SKDelayAddHandler)
	elastic_job.Get().RegisterHandler(model.SKDelayEndTag, svc.SKDelayEndHandler)

	return svc
}

func (s *SKSvc) List(sctx core.SvcContext, pr *page.PageRequest) (*page.Page, error) {
	ctx := sctx.Context()
	mgr := s.SKRepo.Mgr(ctx, s.DB.GetDb(ctx))

	limit, offset := pr.GetLimitAndOffset()
	pr.AddAllowSortField(model.SecKillColumns.CreateTime)
	pr.AddAllowSortField(model.SecKillColumns.StartTime)
	sort, _ := pr.Sort()

	// 默认不筛选出Status异常的SK
	pr.Filter[model.SecKillColumns.Status] = 1

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
	ctx := sctx.Context()
	mgr := s.SKRepo.Mgr(ctx, s.DB.GetDb(ctx))

	// 确保 ShopId 存在
	shopResp, err := s.ShopClient.List(ctx, &proto.ListReq{
		PageNo:    1,
		PageSize:  1,
		SortBy:    "",
		FieldList: []string{"name"},
		ShopIds:   []int64{int64(param.ShopId)},
	})
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	if len(shopResp.GetData()) == 0 {
		return response.NewError(
			http.StatusBadRequest,
			response.ParamBindError,
			"请输入正确的ShopId",
		)
	}

	if time.Until(param.StartTime) < time.Second*30 {
		return response.NewError(
			http.StatusBadRequest,
			response.ParamBindError,
			"start_time 过近",
		)
	}
	if param.EndTime.Sub(param.StartTime) < 0 {
		return response.NewErrorAutoMsg(
			http.StatusBadRequest,
			response.ParamBindError,
		)
	}

	bean := &model.SecKill{
		ShopID:     param.ShopId,
		StartTime:  int(param.StartTime.Unix()),
		EndTime:    int(param.EndTime.Unix()),
		Status:     model.SKWait,
		CreateTime: int(time.Now().Unix()),
	}

	// 存数据库（试用一下事务）
	tx := s.DB.GetDb(ctx).Begin()
	mgr.UpdateDB(tx)

	err = mgr.CreateSecKill(bean)
	if err != nil {
		tx.Rollback()
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}

	if bean.ID == 0 {
		tx.Rollback()
		return response.NewError(
			http.StatusInternalServerError,
			response.ServerError,
			"数据库插入失败",
		)
	}

	key := fmt.Sprintf("%s/%s", model.SKDelayAddTag, encrypt.MD5(uuid.New().String()))
	err = elastic_job.Get().AddJob(&elastic_job.Job{
		Key:       key,
		DelayTime: param.StartTime.Unix(),
		Tag:       model.SKDelayAddTag,
		Args: map[string]interface{}{
			"id": bean.ID,
		},
	})
	if err != nil {
		// 添加任务失败
		tx.Rollback()
		return response.NewError(
			http.StatusInternalServerError,
			response.ServerError,
			"无法开启定时",
		).WithErr(err)
	}

	keyEnd := fmt.Sprintf("%s/%s", model.SKDelayEndTag, encrypt.MD5(uuid.New().String()))
	err = elastic_job.Get().AddJob(&elastic_job.Job{
		Key:       keyEnd,
		DelayTime: param.EndTime.Unix(),
		Tag:       model.SKDelayEndTag,
		Args: map[string]interface{}{
			"id": bean.ID,
		},
	})
	if err != nil {
		tx.Rollback()
		return response.NewError(
			http.StatusInternalServerError,
			response.ServerError,
			"无法开启定时-End",
		).WithErr(err)
	}

	tx.Commit()
	mgr.Reset()
	return nil
}

func (s *SKSvc) Delete(sctx core.SvcContext, id int) error {
	// 确保 ShopId 存在
	ctx := sctx.Context()
	mgr := s.SKRepo.Mgr(ctx, s.DB.GetDb(ctx))

	err := mgr.UpdateSecKill(&model.SecKill{
		ID:     id,
		Status: model.SKClose,
	})
	if err != nil {
		return response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return nil
}

// SKDelayAddHandler 操作秒杀上架
func (s *SKSvc) SKDelayAddHandler(j *elastic_job.Job) error {
	idInter, ok := j.Args["id"]
	if !ok {
		return fmt.Errorf("unfound the id. ")
	}
	id := cast.ToInt(idInter)

	mgr := s.SKRepo.Mgr(context.Background(), s.DB.GetDb(context.Background()))

	// 检查SecKill状态
	sk, err := mgr.WithOptions(mgr.WithID(id)).WithSelects(model.SecKillColumns.ID, model.SecKillColumns.Status).Get()
	if err != nil {
		return fmt.Errorf("err from db: %w ", err)
	}
	if sk.Status != model.SKWait {
		// 只有 SKWait 状态可以
		// 不是则直接退出（无害操作）
		return nil
	}

	err = mgr.UpdateSecKill(&model.SecKill{
		ID:     id,
		Status: model.SKShopping,
	})
	if err != nil {
		return fmt.Errorf("err from db: %w ", err)
	}

	return nil
}

func (s *SKSvc) SKDelayEndHandler(j *elastic_job.Job) error {
	idInter, ok := j.Args["id"]
	if !ok {
		return fmt.Errorf("unfound the id. ")
	}
	id := cast.ToInt(idInter)

	mgr := s.SKRepo.Mgr(context.Background(), s.DB.GetDb(context.Background()))

	// 检查SecKill状态
	sk, err := mgr.WithOptions(mgr.WithID(id)).WithSelects(model.SecKillColumns.ID, model.SecKillColumns.Status).Get()
	if err != nil {
		return fmt.Errorf("err from db: %w ", err)
	}
	if sk.Status != model.SKShopping {
		// 只有 SKShopping 状态可以
		// 不是则直接退出（无害操作）
		return nil
	}

	err = mgr.UpdateSecKill(&model.SecKill{
		ID:     id,
		Status: model.SKFinish,
	})
	if err != nil {
		return fmt.Errorf("err from db: %w ", err)
	}
	return nil
}
