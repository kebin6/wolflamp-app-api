package exchange

import (
	"context"
	"github.com/kebin6/wolflamp-app-api/internal/logic/player"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/zeromicro/go-zero/core/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExchangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeLogic {
	return &ExchangeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *ExchangeLogic) Exchange(req *types.ExchangeReq) (resp *types.ExchangeResp, err error) {

	if !l.svcCtx.Config.WolfLampRpc.Enabled {
		return nil, errorx.NewApiInternalError("common.wolfLampDisable")
	}

	id, err := player.NewPersonInfoLogic(l.ctx, l.svcCtx).GetPlayerId()
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.WolfLampRpc.Exchange(l.ctx, &wolflamp.ExchangeReq{
		PlayerId:   *id,
		Type:       req.Type,
		CoinAmount: req.CoinAmount,
		LampAmount: req.LampAmount,
	})
	return &types.ExchangeResp{
		Data: rpcResp.Id,
	}, err

}
