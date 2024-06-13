package player

import (
	"context"
	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/suyuan32/simple-admin-common/config"
	"github.com/suyuan32/simple-admin-common/utils/encrypt"
	"github.com/zeromicro/go-zero/core/errorx"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.LoginResp, err error) {

	if !l.svcCtx.Config.WolfLampRpc.Enabled {
		return nil, errorx.NewInternalError("common.wolfLampDisable")
	}

	// 校验验证码
	cacheKey := config.RedisCaptchaPrefix + req.Email
	//captchaCode, err := l.svcCtx.Redis.GetDel(l.ctx, cacheKey).Result()
	captchaCode, err := l.svcCtx.Redis.Get(l.ctx, cacheKey).Result()
	if err != nil {
		return nil, errorx.NewCodeInvalidArgumentError("captcha.wrong")
	}

	if captchaCode != req.CaptchaCode {
		return nil, errorx.NewCodeInvalidArgumentError("captcha.wrong")
	}

	info, err := l.svcCtx.WolfLampRpc.GetByEmail(l.ctx, &wolflamp.GetByEmailReq{Email: req.Email})

	if err != nil && status.Convert(err).Message() != "target does not exist" {
		return nil, err
	}

	if info != nil {
		return nil, errorx.NewCodeInvalidArgumentError("common.emailExist")
	}

	//inviter, err := l.svcCtx.WolfLampRpc.GetByInviteCode(l.ctx, &wolflamp.GetByInviteCodeReq{InviteCode: req.InvitedCode})
	//if err != nil {
	//	if status.Convert(err).Message() != i18n.TargetNotFound {
	//		return nil, errorx.NewCodeInvalidArgumentError("common.inviterNotFound")
	//	}
	//	return nil, err
	//}

	_, err = l.svcCtx.WolfLampRpc.CreatePlayer(l.ctx, &wolflamp.CreatePlayerReq{
		Name:                 req.Email,
		Email:                req.Email,
		Password:             encrypt.BcryptEncrypt(req.Password),
		InvitedCode:          req.InvitedCode,
		Status:               1,
		Lamp:                 0,
		Rank:                 0,
		Amount:               0,
		InvitedNum:           0,
		TotalIncome:          0,
		ProfitAndLoss:        0,
		Recent_100WinPercent: 0,
		SystemCommission:     0,
	})
	if err != nil {
		return nil, err
	}

	// 注册完以后直接登陆
	return NewLoginLogic(l.ctx, l.svcCtx).Login(&types.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})

}
