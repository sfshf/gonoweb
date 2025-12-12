// ping_handles包下放置ping等测试型处理函数
package ping

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Ping 测试服务路由是否正常服务
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "PONG %s", time.Now().Format("2006-01-02 15:04:05.000000000"))
}
