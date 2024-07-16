package player

import (
	"context"
	"encoding/json"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/zeromicro/go-zero/core/errorx"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyLogic {
	return &NotifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *NotifyLogic) Notify(req *types.NotifyReq) (resp *types.BaseMsgResp, err error) {

	// 将req转为json字符串
	reqJson, err := json.Marshal(req)
	if err != nil {
		l.Logger.Infof("Notify receive data error: %s", "Invalid Request Parameters")
		return nil, errorx.NewApiBadRequestError("Invalid Request Parameters")
	}

	l.Logger.Infof("Notify Receive Data: %s", reqJson)
	if !l.svcCtx.Config.WolfLampRpc.Enabled {
		return nil, errorx.NewInternalError("common.wolfLampDisable")
	}

	// 校验签名
	validateResp, err := l.svcCtx.WolfLampRpc.ValidateGcicsSign(l.ctx, &wolflamp.ValidateGcicsSignReq{
		Timestamp: req.Time, Sign: req.Sign,
	})
	if err != nil {
		return nil, err
	}
	if !validateResp.IsValid {
		return nil, errorx.NewApiForbiddenError("Invalid Sign")
	}

	id, err := strconv.ParseUint(req.OrderId, 10, 64)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.WolfLampRpc.Notify(l.ctx,
		&wolflamp.NotifyExchangeReq{
			Id:     id,
			IsPaid: req.PaymentStatus == 1,
			Amount: req.Amount,
		})
	if err != nil {
		return nil, err
	}

	return &types.BaseMsgResp{}, nil
}
