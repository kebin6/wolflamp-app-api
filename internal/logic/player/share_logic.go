package player

import (
	"context"
	"fmt"

	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareLogic {
	return &ShareLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}

}

func (l *ShareLogic) Share() (resp *types.ShareResp, err error) {

	//info, err := NewPersonInfoLogic(l.ctx, l.svcCtx).PersonInfo()
	//if err != nil {
	//	return nil, err
	//}

	shareLink := fmt.Sprintf(l.svcCtx.Config.ProjectConf.ShareLink)

	// 生成二维码
	//var qrcodeImage []byte
	//qrcodeImage, err = qrcode.Encode(shareLink, qrcode.Medium, 256)
	//if err != nil {
	//	return nil, err
	//}
	//qrcodeImageBase64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrcodeImage)

	return &types.ShareResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 0,
			Msg:  "success",
		},
		Data: shareLink,
	}, nil

}
