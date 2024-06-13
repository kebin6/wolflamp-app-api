package exchange

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kebin6/wolflamp-app-api/internal/logic/exchange"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"
)

// swagger:route post /exchange exchange Exchange
//
// Exchange Coin To Lamb or Exchange Lamb To Coin | 兑换接口（前置条件：登陆）
//
// Exchange Coin To Lamb or Exchange Lamb To Coin | 兑换接口（前置条件：登陆）
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: ExchangeReq
//
// Responses:
//  200: ExchangeResp

func ExchangeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExchangeReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := exchange.NewExchangeLogic(r.Context(), svcCtx)
		resp, err := l.Exchange(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
