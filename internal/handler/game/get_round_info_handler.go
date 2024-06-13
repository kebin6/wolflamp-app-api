package game

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kebin6/wolflamp-app-api/internal/logic/game"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"
)

// swagger:route post /game game GetRoundInfo
//
// 获取游戏盘口数据接口
//
// 获取游戏盘口数据接口
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: RoundReq
//
// Responses:
//  200: RoundResp

func GetRoundInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RoundReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := game.NewGetRoundInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetRoundInfo(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
