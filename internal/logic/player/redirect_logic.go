package player

import (
	"context"
	"fmt"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/zeromicro/go-zero/core/errorx"
	"google.golang.org/grpc/status"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type RedirectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRedirectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedirectLogic {
	return &RedirectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *RedirectLogic) Redirect(req *types.RedirectReq) (resp *types.RedirectResp, err error) {

	if !l.svcCtx.Config.WolfLampRpc.Enabled {
		return nil, errorx.NewInternalError("common.wolfLampDisable")
	}

	// 校验签名
	validateResp, err := l.svcCtx.WolfLampRpc.ValidateGcicsSign(l.ctx, &wolflamp.ValidateGcicsSignReq{
		Timestamp: req.Time, Sign: req.Sign,
	})
	if err != nil {
		return nil, err
	}
	if !validateResp.IsValid {
		return nil, errorx.NewApiForbiddenError("invalid sign")
	}

	info, err := l.svcCtx.WolfLampRpc.FindPlayer(l.ctx, &wolflamp.FindPlayerReq{GcicsUserId: &req.UserId})

	if err != nil && status.Convert(err).Message() != "target does not exist" {
		return nil, err
	}

	// 用户未创建先创建用户信息
	if err != nil && status.Convert(err).Message() == "target does not exist" {
		createResp, err := l.svcCtx.WolfLampRpc.CreatePlayer(l.ctx, &wolflamp.CreatePlayerReq{
			Name:             "CuteLamb",
			Email:            strconv.FormatUint(req.UserId, 10),
			Password:         "",
			InvitedCode:      "solamb0",
			Status:           1,
			CoinLamb:         0,
			TokenLamb:        0,
			Rank:             0,
			Amount:           0,
			InvitedNum:       0,
			TotalIncome:      0,
			ProfitAndLoss:    0,
			SystemCommission: 0,
			GcicsUserId:      req.UserId,
			GcicsToken:       req.Token,
			ReturnUrl:        req.ReturnUrl,
		})
		if err != nil {
			return nil, err
		}

		info = &wolflamp.PlayerInfo{Id: createResp.Id}
	}

	// 注册完以后直接登陆
	token, err := NewLoginLogic(l.ctx, l.svcCtx).GenerateToken(info)
	if err != nil {
		return nil, err
	}

	return &types.RedirectResp{
		Data: types.RedirectInfo{
			GameUrl: fmt.Sprintf(l.svcCtx.Config.ProjectConf.RedirectGameUrl, *token),
		},
	}, nil
}
