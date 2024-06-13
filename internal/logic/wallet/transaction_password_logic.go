package wallet

import (
	"context"
	"github.com/kebin6/wolflamp-app-api/internal/logic/player"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/suyuan32/simple-admin-common/utils/encrypt"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransactionPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransactionPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransactionPasswordLogic {
	return &TransactionPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *TransactionPasswordLogic) TransactionPassword(req *types.TransactionPasswordReq) (resp *types.BaseMsgResp, err error) {

	_, err = NewVerifyTransactionPasswordLogic(l.ctx, l.svcCtx).
		VerifyTransactionPassword(&types.VerifyTransactionPasswordReq{OldPassword: req.OldPassword})
	if err != nil {
		return nil, err
	}

	id, err := player.NewPersonInfoLogic(l.ctx, l.svcCtx).GetPlayerId()
	if err != nil {
		return nil, err
	}

	encryptNewPassword := encrypt.BcryptEncrypt(req.NewPassword)
	_, err = l.svcCtx.WolfLampRpc.UpdatePlayer(l.ctx, &wolflamp.UpdatePlayerReq{
		Id:                  *id,
		TransactionPassword: &encryptNewPassword,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseMsgResp{}, nil

}
