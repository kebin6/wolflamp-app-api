package reward

import (
	"context"
	"github.com/kebin6/wolflamp-app-api/internal/logic/player"
	"github.com/kebin6/wolflamp-rpc/common/enum/rewardenum"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/zeromicro/go-zero/core/errorx"

	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindLogic {
	return &FindLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *FindLogic) Find(req *types.FindRewardReq) (resp *types.FindRewardResp, err error) {

	if !l.svcCtx.Config.WolfLampRpc.Enabled {
		return nil, errorx.NewCodeUnavailableError(i18n.ServiceUnavailable)
	}

	id, err := player.NewPersonInfoLogic(l.ctx, l.svcCtx).GetPlayerId()
	if err != nil {
		return nil, err
	}

	data, err := l.svcCtx.WolfLampRpc.FindReward(l.ctx, &wolflamp.FindRewardReq{Id: req.Id, ToId: id})
	if err != nil {
		return nil, err
	}

	return &types.FindRewardResp{
		Data: l.Po2Vo(data)}, nil

}

func (l *FindLogic) Po2Vo(po *wolflamp.RewardInfo) (vo types.RewardInfo) {

	return types.RewardInfo{
		Id:               po.Id,
		ToId:             po.ToId,
		Remark:           po.Remark,
		Num:              po.Num,
		CreateAt:         po.CreatedAt,
		Status:           po.Status,
		StatusDesc:       rewardenum.New(po.Status).Desc(),
		ContributorId:    po.ContributorId,
		ContributorEmail: po.ContributorEmail,
		ContributorLevel: po.ContributorLevel,
	}

}
