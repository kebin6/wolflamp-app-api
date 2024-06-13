package player

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kebin6/wolflamp-app-api/internal/logic/player"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"
)

// swagger:route post /login player Login
//
// 登陆接口
//
// 登陆接口
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: LoginReq
//
// Responses:
//  200: LoginResp

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := player.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
