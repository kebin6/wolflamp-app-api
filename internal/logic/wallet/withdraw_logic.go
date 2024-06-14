package wallet

import (
	"context"
	"fmt"
	"github.com/kebin6/wolflamp-app-api/internal/logic/player"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"
	"github.com/kebin6/wolflamp-rpc/common/enum/orderenum"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/errorx"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type WithdrawLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawLogic {
	return &WithdrawLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *WithdrawLogic) Withdraw(req *types.WithdrawReq) (resp *types.WithdrawResp, err error) {

	if !l.svcCtx.Config.WolfLampRpc.Enabled {
		return nil, errorx.NewCodeUnavailableError(i18n.ServiceUnavailable)
	}

	_, err = NewVerifyTransactionPasswordLogic(l.ctx, l.svcCtx).
		VerifyTransactionPassword(&types.VerifyTransactionPasswordReq{OldPassword: req.Password})
	if err != nil {
		return nil, err
	}

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

	// 检查是否存在审核中的提币订单
	orders, err := l.svcCtx.WolfLampRpc.ListOrder(l.ctx, &wolflamp.ListOrderReq{
		PlayerId: &info.Id, Type: pointy.GetPointer(orderenum.Withdraw.Val()),
		Status: pointy.GetPointer(orderenum.Applying.Val()), PageSize: 1, Page: 1})
	if err != nil {
		return nil, err
	}
	if len(orders.Data) > 0 {
		return nil, errorx.NewCodeAlreadyExistsError("wallet.withdrawExist")
	}

	// 根据系统配置校验提币数量是否符合要求
	minWithdrawNumConfig, err := l.svcCtx.WolfLampRpc.GetMinWithdrawNum(l.ctx, &wolflamp.Empty{})
	if err != nil {
		return nil, err
	}
	if req.Amount <= float64(minWithdrawNumConfig.MinWithdrawNum) {
		return nil, errorx.NewCodeInvalidArgumentError(
			fmt.Sprintf(l.svcCtx.Trans.Trans(l.ctx, "wallet.minWithdrawLimit"), minWithdrawNumConfig.MinWithdrawNum))
	}

	// 判断是否需要进入审核状态
	withdrawThreshold, err := l.svcCtx.WolfLampRpc.GetWithdrawThreshold(l.ctx, &wolflamp.Empty{})
	if err != nil {
		return nil, err
	}
	statusVal := orderenum.Pending.Val()
	if req.Amount >= float64(withdrawThreshold.WithdrawThreshold) {
		statusVal = orderenum.Applying.Val()
	}

	// 计算手续费
	// 提币订单需要手续费
	calRes, err := l.svcCtx.WolfLampRpc.
		CalculateWithdrawFee(l.ctx, &wolflamp.CalculateWithdrawFeeReq{Balance: info.Amount, Amount: req.Amount})
	if err != nil {
		return nil, err
	}
	handlingFee := calRes.HandlingFee

	orderCreateResp, err := l.svcCtx.WolfLampRpc.CreateOrder(l.ctx, &wolflamp.CreateOrderReq{
		Status: statusVal, Type: orderenum.Withdraw.Val(), Num: req.Amount, ToAddress: req.WithdrawAddress,
		FromAddress: info.DepositAddress, PlayerId: *id, HandlingFee: handlingFee,
		Network: l.svcCtx.Config.ProjectConf.Network})
	if err != nil {
		return nil, err
	}

	return &types.WithdrawResp{
		Data: orderCreateResp.Id,
	}, nil

}
