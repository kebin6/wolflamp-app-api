package reward

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

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *ListLogic) List(req *types.ListRewardReq) (resp *types.ListRewardResp, err error) {

	if !l.svcCtx.Config.WolfLampRpc.Enabled {
		return nil, errorx.NewCodeUnavailableError(i18n.ServiceUnavailable)
	}

	id, err := player.NewPersonInfoLogic(l.ctx, l.svcCtx).GetPlayerId()
	if err != nil {
		return nil, err
	}

	data, err := l.svcCtx.WolfLampRpc.ListReward(l.ctx,
		&wolflamp.ListRewardReq{
			Page:             req.Page,
			PageSize:         req.PageSize,
			ToId:             id,
			ContributorLevel: &req.ContributorLevel,
		})
	if err != nil {
		return nil, err
	}
	resp = &types.ListRewardResp{}
	resp.Msg = l.svcCtx.Trans.Trans(l.ctx, i18n.Success)
	resp.Data.Total = data.GetTotal()

	findLogic := NewFindLogic(l.ctx, l.svcCtx)
	for _, v := range data.Data {
		resp.Data.Data = append(resp.Data.Data, findLogic.Po2Vo(v))
	}
	return resp, nil

}
