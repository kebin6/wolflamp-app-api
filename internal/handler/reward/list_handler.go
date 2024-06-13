package reward

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kebin6/wolflamp-app-api/internal/logic/reward"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"
)

// swagger:route post /reward/list reward List
//
// Get Reward Log List | 奖励列表（前置条件：登陆）
//
// Get Reward Log List | 奖励列表（前置条件：登陆）
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: ListRewardReq
//
// Responses:
//  200: ListRewardResp

func ListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListRewardReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := reward.NewListLogic(r.Context(), svcCtx)
		resp, err := l.List(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
