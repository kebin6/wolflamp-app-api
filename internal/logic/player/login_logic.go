package player

import (
	"context"
	"fmt"
	"github.com/kebin6/wolflamp-rpc/common/enum/cachekey"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/suyuan32/simple-admin-common/utils/encrypt"
	"github.com/suyuan32/simple-admin-common/utils/jwt"
	"github.com/zeromicro/go-zero/core/errorx"
	"google.golang.org/grpc/status"
	"time"

	"github.com/kebin6/wolflamp-app-api/internal/svc"
	"github.com/kebin6/wolflamp-app-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {

	if !l.svcCtx.Config.WolfLampRpc.Enabled {
		return nil, errorx.NewInternalError("common.wolfLampDisable")
	}

	info, err := l.svcCtx.WolfLampRpc.GetByEmail(l.ctx, &wolflamp.GetByEmailReq{Email: req.Email})

	if err != nil {
		if status.Convert(err).Message() == "target does not exist" {
			return nil, errorx.NewCodeInvalidArgumentError("common.playerNotFound")
		}
		return nil, err
	}

	// 验证密码是否正确
	if ok := encrypt.BcryptCheck(req.Password, info.Password); !ok {
		return nil, errorx.NewCodeInvalidArgumentError("common.wrongPassword")
	}

	token, err := l.GenerateToken(info)
	return &types.LoginResp{
		Data: types.LoginInfo{
			Info: types.PlayerInfo{
				Id: info.Id,
			},
			Token: *token,
		},
	}, nil

}

// GenerateToken 生成登陆token
func (l *LoginLogic) GenerateToken(info *wolflamp.PlayerInfo) (*string, error) {
	token, err := jwt.NewJwtToken(
		l.svcCtx.Config.Auth.AccessSecret,
		time.Now().Unix(),
		l.svcCtx.Config.Auth.AccessExpire,
		jwt.WithOption("id", info.Id))
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf(string(cachekey.GameAuthToken), info.Id)
	// 先删除旧token
	err = l.svcCtx.Redis.Del(l.ctx, cacheKey).Err()
	if err != nil {
		logx.Errorw("failed to delete token in redis", logx.Field("detail", err))
		return nil, err
	}

	// 缓存现有token
	err = l.svcCtx.Redis.Set(l.ctx, cacheKey, token, time.Duration(l.svcCtx.Config.Auth.AccessExpire)).Err()
	if err != nil {
		logx.Errorw("failed to set token in redis", logx.Field("detail", err))
		return nil, err
	}
	return &token, nil
}
