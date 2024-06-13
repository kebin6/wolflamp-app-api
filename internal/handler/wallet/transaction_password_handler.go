package wallet

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kebin6/wolflamp-app-api/internal/logic/wallet"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"
)

// swagger:route post /transaction_password/update wallet TransactionPassword
//
// 设置交易密码接口（前置条件：登陆）
//
// 设置交易密码接口（前置条件：登陆）
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: TransactionPasswordReq
//
// Responses:
//  200: BaseMsgResp

func TransactionPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TransactionPasswordReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := wallet.NewTransactionPasswordLogic(r.Context(), svcCtx)
		resp, err := l.TransactionPassword(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
