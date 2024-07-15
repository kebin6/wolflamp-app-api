package exchange

import (
	"context"
	"github.com/kebin6/wolflamp-app-api/internal/logic/player"
	"github.com/kebin6/wolflamp-rpc/common/enum/exchangeenum"
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

func (l *FindLogic) Find(req *types.FindExchangeReq) (resp *types.FindExchangeResp, err error) {

	if !l.svcCtx.Config.WolfLampRpc.Enabled {
		return nil, errorx.NewCodeUnavailableError(i18n.ServiceUnavailable)
	}

	id, err := player.NewPersonInfoLogic(l.ctx, l.svcCtx).GetPlayerId()
	if err != nil {
		return nil, err
	}

	data, err := l.svcCtx.WolfLampRpc.FindExchange(l.ctx, &wolflamp.FindExchangeReq{Id: req.Id, PlayerId: id})
	if err != nil {
		return nil, err
	}

	return &types.FindExchangeResp{
		Data: l.Po2Vo(data),
	}, nil

}

func (l *FindLogic) Po2Vo(po *wolflamp.ExchangeInfo) (vo types.ExchangeInfo) {

	return types.ExchangeInfo{
		Id:            po.Id,
		CreateAt:      po.CreatedAt,
		Status:        po.Status,
		StatusDesc:    exchangeenum.NewExchangeStatus(po.Status).Desc(),
		TransactionId: po.TransactionId,
		Type:          po.Type,
		TypeDesc:      exchangeenum.NewExchangeType(po.Type).Desc(),
		CoinAmount:    po.CoinNum,
		LampAmount:    po.LampNum,
		Mode:          po.Mode,
	}

}
