package order

import (
	"context"
	"github.com/kebin6/wolflamp-app-api/internal/logic/player"
	"github.com/kebin6/wolflamp-rpc/common/enum/orderenum"
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

func (l *FindLogic) Find(req *types.FindOrderReq) (resp *types.FindOrderResp, err error) {

	if !l.svcCtx.Config.WolfLampRpc.Enabled {
		return nil, errorx.NewCodeUnavailableError(i18n.ServiceUnavailable)
	}

	id, err := player.NewPersonInfoLogic(l.ctx, l.svcCtx).GetPlayerId()
	if err != nil {
		return nil, err
	}

	data, err := l.svcCtx.WolfLampRpc.FindOrder(l.ctx, &wolflamp.FindOrderReq{Id: req.Id, PlayerId: id})
	if err != nil {
		return nil, err
	}

	return &types.FindOrderResp{
		Data: l.Po2Vo(data),
	}, nil

}

func (l *FindLogic) Po2Vo(po *wolflamp.OrderInfo) (vo types.OrderInfo) {

	return types.OrderInfo{
		Id:            po.Id,
		TransactionId: po.TransactionId,
		Type:          po.Type,
		Num:           po.Num,
		CreateAt:      po.CreatedAt,
		Status:        po.Status,
		StatusDesc:    orderenum.NewOrderStatus(po.Status).Desc(),
		ToAddress:     po.ToAddress,
		Network:       po.Network,
		HandlingFee:   po.HandlingFee,
	}

}
