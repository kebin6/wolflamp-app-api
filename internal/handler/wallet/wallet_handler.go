package wallet

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kebin6/wolflamp-app-api/internal/logic/wallet"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
)

// swagger:route post /wallet wallet Wallet
//
// 获取钱包详情接口（前置条件：登陆）
//
// 获取钱包详情接口（前置条件：登陆）
//
// Responses:
//  200: WalletResp

func WalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := wallet.NewWalletLogic(r.Context(), svcCtx)
		resp, err := l.Wallet()
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
