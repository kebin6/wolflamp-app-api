package player

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kebin6/wolflamp-app-api/internal/logic/player"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
)

// swagger:route post /share player Share
//
// 获取分享链接接口（前置条件：登陆）
//
// 获取分享链接接口（前置条件：登陆）
//
// Responses:
//  200: ShareResp

func ShareHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := player.NewShareLogic(r.Context(), svcCtx)
		resp, err := l.Share()
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
