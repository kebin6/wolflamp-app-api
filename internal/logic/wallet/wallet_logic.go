package wallet

import (
	"context"
	"github.com/kebin6/wolflamp-app-api/internal/logic/player"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/zeromicro/go-zero/core/errorx"
	"google.golang.org/grpc/status"
	"net/http"

	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WalletLogic {
	return &WalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *WalletLogic) Wallet() (resp *types.WalletResp, err error) {

	if !l.svcCtx.Config.WolfLampRpc.Enabled {
		return nil, errorx.NewApiInternalError("common.wolfLampDisable")
	}

	id, err := player.NewPersonInfoLogic(l.ctx, l.svcCtx).GetPlayerId()
	if err != nil {
		return nil, err
	}

	info, err := l.svcCtx.WolfLampRpc.FindPlayer(l.ctx, &wolflamp.FindPlayerReq{Id: id})

	if err != nil {
		if status.Convert(err).Message() != i18n.TargetNotFound {
			return nil, errorx.NewApiInternalError("common.playerNotFound")
		}
		return nil, err
	}

	balance, err := l.svcCtx.WolfLampRpc.GetGcicsBalance(l.ctx, &wolflamp.GetGcicsBalanceReq{PlayerId: *id})
	if err != nil {
		if status.Convert(err).Message() == "Unauthorized" {
			return nil, errorx.NewApiErrorWithoutMsg(http.StatusUnauthorized)
		}
		return nil, err
	}

	isInit := uint32(2)
	if info.TransactionPassword != "" {
		isInit = 1
	}
	return &types.WalletResp{
		Data: types.WalletInfo{
			IsInit:       isInit,
			CoinBalance:  balance.Coin,
			TokenBalance: balance.Token,
			CoinLamb:     info.CoinLamb,
			TokenLamb:    info.TokenLamb,
		},
	}, nil
}
