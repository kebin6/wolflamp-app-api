package player

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kebin6/wolflamp-app-api/internal/logic/player"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"
)

// swagger:route post /change player ForgetPassword
//
// 忘记密码接口（修改完密码以后会返回新的登陆信息）
//
// 忘记密码接口（修改完密码以后会返回新的登陆信息）
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: ChangePasswordReq
//
// Responses:
//  200: LoginResp

func ForgetPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChangePasswordReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := player.NewForgetPasswordLogic(r.Context(), svcCtx)
		resp, err := l.ForgetPassword(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
