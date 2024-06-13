package player

import (
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/random"
	"github.com/suyuan32/simple-admin-common/config"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/suyuan32/simple-admin-message-center/types/mcms"
	"github.com/zeromicro/go-zero/core/errorx"
	"strconv"
	"time"

	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaLogic {
	return &CaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *CaptchaLogic) Captcha(req *types.CaptchaReq) (resp *types.BaseMsgResp, err error) {

	if !l.svcCtx.Config.McmsRpc.Enabled {
		return nil, errorx.NewCodeInvalidArgumentError("captcha.mcmsNotEnabled")
	}

	cacheKey := config.RedisCaptchaPrefix + req.Email
	captcha := random.RandInt(100000, 999999)
	// 验证码过期时间/秒
	captchaExpireTime := time.Duration(l.svcCtx.Config.ProjectConf.EmailCaptchaExpiredTime) * time.Second
	_, err = l.svcCtx.McmsRpc.SendEmail(l.ctx, &mcms.EmailInfo{
		Target:  []string{req.Email},
		Subject: fmt.Sprintf(l.svcCtx.Trans.Trans(l.ctx, "captcha.email.subject"), l.svcCtx.Trans.Trans(l.ctx, "system.name")),
		Content: fmt.Sprintf(l.svcCtx.Trans.Trans(l.ctx, "captcha.email.template"), captcha, captchaExpireTime/time.Minute),
	})
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.Redis.Set(l.ctx, cacheKey, strconv.Itoa(captcha), captchaExpireTime).Err()
	if err != nil {
		logx.Errorw("failed to write email captcha to redis", logx.Field("detail", err))
		return nil, errorx.NewCodeInternalError(i18n.RedisError)
	}

	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, i18n.Success)}, nil

}
