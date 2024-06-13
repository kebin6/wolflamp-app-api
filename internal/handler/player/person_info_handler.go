package player

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kebin6/wolflamp-app-api/internal/logic/player"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
)

// swagger:route post /person player PersonInfo
//
// 获取个人信息接口（前置条件：登陆）
//
// 获取个人信息接口（前置条件：登陆）
//
// Responses:
//  200: PersonInfoResp

func PersonInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := player.NewPersonInfoLogic(r.Context(), svcCtx)
		resp, err := l.PersonInfo()
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
