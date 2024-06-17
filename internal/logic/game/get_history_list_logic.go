package game

import (
	"context"
	"github.com/kebin6/wolflamp-app-api/internal/logic/player"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/zeromicro/go-zero/core/errorx"

	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHistoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHistoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHistoryListLogic {
	return &GetHistoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *GetHistoryListLogic) GetHistoryList(req *types.HistoryListReq) (resp *types.HistoryListResp, err error) {

	if !l.svcCtx.Config.WolfLampRpc.Enabled {
		return nil, errorx.NewCodeUnavailableError(i18n.ServiceUnavailable)
	}

	id, err := player.NewPersonInfoLogic(l.ctx, l.svcCtx).GetPlayerId()
	if err != nil {
		return nil, err
	}

	data, err := l.svcCtx.WolfLampRpc.ListHistoryInvest(l.ctx,
		&wolflamp.ListHistoryInvestReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			PlayerId: id,
		})
	if err != nil {
		return nil, err
	}
	resp = &types.HistoryListResp{}
	resp.Msg = l.svcCtx.Trans.Trans(l.ctx, i18n.Success)
	resp.Data.Total = data.GetTotal()

	for _, v := range data.Data {
		resp.Data.Data = append(resp.Data.Data, types.InvestRecordInfo{
			ProfitAndLoss: v.ProfitAndLoss,
			RecordTime:    v.CreatedAt,
		})
	}
	return resp, nil

}
