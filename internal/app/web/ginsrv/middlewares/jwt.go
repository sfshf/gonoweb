package ginmw

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	gono_web "github.com/sfshf/gonoweb/internal/app/web"
	"github.com/sfshf/gonoweb/internal/config"
	user_svc "github.com/sfshf/gonoweb/internal/service/user"
	jwt_util "github.com/sfshf/gonoweb/internal/util/jwt"
)

const (
	GinContext_JWTClaims = "jwt_claims"
)

// Jwt middleware
func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 检查Authorization头部
		jwtString := c.GetHeader("Authorization")
		if jwtString == "" {
			c.JSON(http.StatusUnauthorized, &gono_web.Response{
				Code: gono_web.ResponseCode_RequestError,
				Msg:  "Authorization头部为空",
			})
			c.Abort()
			return
		}
		// 2. 检查用户的token是否是本人当前所在IP上使用的
		if err := user_svc.CheckTokenWithIP(jwtString, c.ClientIP()); err != nil {
			c.JSON(http.StatusUnauthorized, &gono_web.Response{
				Code: gono_web.ResponseCode_RequestError,
				Msg:  err.Error(),
			})
			c.Abort()
			return
		}
		claims, err := jwt_util.ParseToken(jwt_util.DefaultSigningMethod, config.AppConfig.Gin.Jwt.SigningKey, jwtString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, &gono_web.Response{
				Code: gono_web.ResponseCode_RequestError,
				Msg:  err.Error(),
			})
			c.Abort()
			return
		}
		// 3. 检查jwt是否过期，如果过期，则提醒用户重新登录，后台不做续期行为
		if claims.ExpiresAt.Add(-config.AppConfig.Gin.Jwt.Expired * time.Second).Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, &gono_web.Response{
				Code: gono_web.ResponseCode_JwtExpired,
				Msg:  "登录token过期，请重新登录",
			})
			c.Abort()
			return
		}
		// 4. 将用户信息（user_xid, domain_xid, role_xid）设置到请求上下文
		c.Set(GinContext_JWTClaims, claims)
		c.Next()
	}
}

func JwtClaims(c *gin.Context) *jwt_util.JwtClaims {
	claims, exists := c.Get(GinContext_JWTClaims)
	if exists {
		obj, ok := claims.(*jwt_util.JwtClaims)
		if ok {
			return obj
		}
	}
	return nil
}
