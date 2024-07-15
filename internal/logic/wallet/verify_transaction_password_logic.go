package wallet

import (
	"context"
	"github.com/kebin6/wolflamp-app-api/internal/logic/player"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/suyuan32/simple-admin-common/utils/encrypt"
	"github.com/zeromicro/go-zero/core/errorx"
	"google.golang.org/grpc/status"

	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyTransactionPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyTransactionPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyTransactionPasswordLogic {
	return &VerifyTransactionPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *VerifyTransactionPasswordLogic) VerifyTransactionPassword(req *types.VerifyTransactionPasswordReq) (resp *types.BaseMsgResp, err error) {

	if !l.svcCtx.Config.WolfLampRpc.Enabled {
		return nil, errorx.NewApiInternalError("common.wolfLampDisable")
	}

	id, err := player.NewPersonInfoLogic(l.ctx, l.svcCtx).GetPlayerId()
	if err != nil {
		return nil, err
	}
	info, err := l.svcCtx.WolfLampRpc.FindPlayer(l.ctx, &wolflamp.FindPlayerReq{Id: id})

	if err != nil {
		if status.Convert(err).Message() == "target does not exist" {
			return nil, errorx.NewApiInternalError("common.playerNotFound")
		}
		return nil, err
	}

	// 首次设置交易密码，不需要验证旧密码，否则需要校验旧密码是否正确
	if info.TransactionPassword != "" {
		if ok := encrypt.BcryptCheck(req.OldPassword, info.TransactionPassword); !ok {
			return nil, errorx.NewCodeInvalidArgumentError("wallet.wrongOldTransactionPassword")
		}
	}

	return &types.BaseMsgResp{}, nil

}
