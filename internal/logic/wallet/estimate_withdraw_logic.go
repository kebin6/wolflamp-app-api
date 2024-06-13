package wallet

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

type EstimateWithdrawLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEstimateWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EstimateWithdrawLogic {
	return &EstimateWithdrawLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *EstimateWithdrawLogic) EstimateWithdraw(req *types.EstimateWithdrawReq) (resp *types.EstimateWithdrawResp, err error) {

	id, err := player.NewPersonInfoLogic(l.ctx, l.svcCtx).GetPlayerId()
	if err != nil {
		return nil, err
	}

	// 获取玩家信息
	info, err := l.svcCtx.WolfLampRpc.FindPlayer(l.ctx, &wolflamp.FindPlayerReq{Id: *id})
	if err != nil {
		if status.Convert(err).Message() != i18n.TargetNotFound {
			return nil, errorx.NewApiInternalError("common.playerNotFound")
		}
		return nil, err
	}

	// 计算手续费
	// 提币订单需要手续费
	calRes, err := l.svcCtx.WolfLampRpc.
		CalculateWithdrawFee(l.ctx, &wolflamp.CalculateWithdrawFeeReq{Balance: info.Amount, Amount: req.Amount})
	if err != nil {
		return nil, err
	}

	return &types.EstimateWithdrawResp{
		BaseDataInfo: types.BaseDataInfo{Msg: i18n.Success},
		Data: types.EstimateWithdrawDataInfo{
			Amount:         req.Amount,
			TotalAmount:    calRes.TotalAmount,
			ReceivedAmount: calRes.ReceivedAmount,
			HandingFee:     calRes.HandlingFee,
		},
	}, nil

}
