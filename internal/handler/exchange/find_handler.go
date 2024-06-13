package exchange

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kebin6/wolflamp-app-api/internal/logic/exchange"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"
)

// swagger:route post /exchange/find exchange Find
//
// Get Exchange Log Info | 兑换记录详情（前置条件：登陆）
//
// Get Exchange Log Info | 兑换记录详情（前置条件：登陆）
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: FindExchangeReq
//
// Responses:
//  200: FindExchangeResp

func FindHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FindExchangeReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := exchange.NewFindLogic(r.Context(), svcCtx)
		resp, err := l.Find(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
