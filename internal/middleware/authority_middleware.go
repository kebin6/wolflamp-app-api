package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/kebin6/wolflamp-rpc/common/enum/cachekey"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"

	"github.com/suyuan32/simple-admin-common/i18n"
)

type AuthorityMiddleware struct {
	Rds   redis.UniversalClient
	Trans *i18n.Translator
}

func NewAuthorityMiddleware(cbn *casbin.Enforcer, rds redis.UniversalClient, trans *i18n.Translator) *AuthorityMiddleware {
	return &AuthorityMiddleware{
		Rds:   rds,
		Trans: trans,
	}
}

func (m *AuthorityMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value("id")

		// check jwt blacklist
		_, err := m.Rds.Get(context.Background(), fmt.Sprintf(string(cachekey.GameAuthToken), id)).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			httpx.Error(w, errorx.NewApiError(http.StatusInternalServerError, err.Error()))
			return
		}
		if err != nil && errors.Is(err, redis.Nil) {
			httpx.Error(w, errorx.NewApiErrorWithoutMsg(http.StatusUnauthorized))
			return
		}
	}
}
