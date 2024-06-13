package player

import (
	"context"
	"encoding/json"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/zeromicro/go-zero/core/errorx"
	"google.golang.org/grpc/status"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type PersonInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPersonInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PersonInfoLogic {
	return &PersonInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *PersonInfoLogic) PersonInfo() (resp *types.PersonInfoResp, err error) {

	if !l.svcCtx.Config.WolfLampRpc.Enabled {
		return nil, errorx.NewApiInternalError("common.wolfLampDisable")
	}

	id, err := l.GetPlayerId()
	if err != nil {
		return nil, err
	}
	info, err := l.svcCtx.WolfLampRpc.FindPlayer(l.ctx, &wolflamp.FindPlayerReq{Id: *id})

	if err != nil {
		if status.Convert(err).Message() != "target does not exist" {
			return nil, errorx.NewApiInternalError("common.playerNotFound")
		}
		return nil, err
	}

	return &types.PersonInfoResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 0,
			Msg:  "success",
		},
		Data: types.PlayerInfo{
			Id:         info.Id,
			Email:      info.Email,
			InviteCode: info.InviteCode,
			Amount:     info.Amount,
			Lamp:       info.Lamp,
		},
	}, nil

}

func (l *PersonInfoLogic) GetPlayerId() (*uint64, error) {
	idNumber, ok := l.ctx.Value("id").(json.Number)
	if !ok {
		return nil, errorx.NewApiInternalError("common.invalidToken")
	}
	id, err := strconv.ParseUint(idNumber.String(), 10, 64)
	if err != nil {
		return nil, errorx.NewApiInternalError("common.invalidToken")
	}
	return &id, err
}
