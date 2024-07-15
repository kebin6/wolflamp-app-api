package notice

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kebin6/wolflamp-app-api/internal/logic/notice"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
)

// swagger:route post /notice notice Notice
//
// 获取最新公告接口
//
// 获取最新公告接口
//
// Responses:
//  200: BaseDataInfo

func NoticeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notice.NewNoticeLogic(r.Context(), svcCtx)
		resp, err := l.Notice()
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
