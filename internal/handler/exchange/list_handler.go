package exchange

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kebin6/wolflamp-app-api/internal/logic/exchange"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"
)

// swagger:route post /exchange/list exchange List
//
// Get Exchange Log List | 兑换记录列表（前置条件：登陆）
//
// Get Exchange Log List | 兑换记录列表（前置条件：登陆）
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: ListExchangeReq
//
// Responses:
//  200: ListExchangeResp

func ListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListExchangeReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := exchange.NewListLogic(r.Context(), svcCtx)
		resp, err := l.List(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
