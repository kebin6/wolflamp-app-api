package game

import (
	"context"
	"github.com/kebin6/wolflamp-app-api/internal/logic/player"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/zeromicro/go-zero/core/errorx"
	"google.golang.org/grpc/status"

	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InvestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInvestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InvestLogic {
	return &InvestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *InvestLogic) Invest(req *types.InvestReq) (resp *types.BaseDataInfo, err error) {

	if !l.svcCtx.Config.WolfLampRpc.Enabled {
		return nil, errorx.NewCodeUnavailableError(i18n.ServiceUnavailable)
	}

	id, err := player.NewPersonInfoLogic(l.ctx, l.svcCtx).GetPlayerId()
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.WolfLampRpc.Invest(l.ctx,
		&wolflamp.CreateInvestReq{RoundId: req.RoundId, PlayerId: *id, FoldNo: req.FoldNo, LambNum: req.LambNum})
	if err != nil {
		if status.Convert(err).Message() == "target does not exist" {
			return nil, errorx.NewCodeInvalidArgumentError("game.roundNotFound")
		}
		return nil, err
	}

	return &types.BaseDataInfo{Code: 0, Msg: i18n.Success}, nil

}
