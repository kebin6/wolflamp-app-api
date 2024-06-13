package wallet

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kebin6/wolflamp-app-api/internal/logic/wallet"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"
)

// swagger:route post /transaction_password/verify wallet VerifyTransactionPassword
//
// 校验交易密码是否正确接口（前置条件：登陆）
//
// 校验交易密码是否正确接口（前置条件：登陆）
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: VerifyTransactionPasswordReq
//
// Responses:
//  200: BaseMsgResp

func VerifyTransactionPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VerifyTransactionPasswordReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := wallet.NewVerifyTransactionPasswordLogic(r.Context(), svcCtx)
		resp, err := l.VerifyTransactionPassword(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
