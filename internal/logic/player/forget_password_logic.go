package player

import (
	"context"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/suyuan32/simple-admin-common/config"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/suyuan32/simple-admin-common/utils/encrypt"
	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/errorx"
	"google.golang.org/grpc/status"

	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ForgetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewForgetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForgetPasswordLogic {
	return &ForgetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *ForgetPasswordLogic) ForgetPassword(req *types.ChangePasswordReq) (resp *types.LoginResp, err error) {

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

	if err != nil {
		if status.Convert(err).Message() != i18n.TargetNotFound {
			return nil, errorx.NewCodeInvalidArgumentError("common.playerNotFound")
		}
		return nil, err
	}

	// 验证密码是否正确
	if ok := encrypt.BcryptCheck(req.OldPassword, info.Password); !ok {
		return nil, errorx.NewCodeInvalidArgumentError("common.wrongPassword")
	}

	// 修改密码
	_, err = l.svcCtx.WolfLampRpc.UpdatePlayer(l.ctx, &wolflamp.UpdatePlayerReq{
		Id:       info.Id,
		Password: pointy.GetPointer(encrypt.BcryptEncrypt(req.NewPassword)),
	})
	if err != nil {
		return nil, err
	}

	// 修改完密码以后直接登陆
	return NewLoginLogic(l.ctx, l.svcCtx).Login(&types.LoginReq{
		Email:    req.Email,
		Password: req.NewPassword,
	})

}
