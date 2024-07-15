package player

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kebin6/wolflamp-app-api/internal/logic/player"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"
)

// swagger:route post /notify player Notify
//
// 提供给GCICS平台的支付回调接口
//
// 提供给GCICS平台的支付回调接口
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: NotifyReq
//
// Responses:
//  200: BaseMsgResp

func NotifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NotifyReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := player.NewNotifyLogic(r.Context(), svcCtx)
		_, err := l.Notify(&req)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write([]byte("SUCCESS"))
	}
}
