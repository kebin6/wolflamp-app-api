package svc

import (
	"github.com/kebin6/wolflamp-app-api/internal/config"
	i18n2 "github.com/kebin6/wolflamp-app-api/internal/i18n"
	"github.com/kebin6/wolflamp-app-api/internal/middleware"
	"github.com/kebin6/wolflamp-rpc/wolflampclient"
	"github.com/redis/go-redis/v9"
	"github.com/suyuan32/simple-admin-message-center/mcmsclient"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/suyuan32/simple-admin-common/i18n"

	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config      config.Config
	Casbin      *casbin.Enforcer
	Authority   rest.Middleware
	Trans       *i18n.Translator
	McmsRpc     mcmsclient.Mcms
	Redis       redis.UniversalClient
	WolfLampRpc wolflampclient.Wolflamp
}

func NewServiceContext(c config.Config) *ServiceContext {

	rds := c.RedisConf.MustNewUniversalRedis()

	cbn := c.CasbinConf.MustNewCasbinWithOriginalRedisWatcher(c.CasbinDatabaseConf.Type, c.CasbinDatabaseConf.GetDSN(), c.RedisConf)

	trans := i18n.NewTranslator(i18n2.LocaleFS)

	return &ServiceContext{
		Config:      c,
		Authority:   middleware.NewAuthorityMiddleware(cbn, rds, trans).Handle,
		Trans:       trans,
		Redis:       rds,
		McmsRpc:     mcmsclient.NewMcms(zrpc.NewClientIfEnable(c.McmsRpc)),
		WolfLampRpc: wolflampclient.NewWolflamp(zrpc.MustNewClient(c.WolfLampRpc)),
	}
}
